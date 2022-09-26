// Copyright (c) 2022 Braden Nicholson

package atlas

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
	"udap/internal/log"
)

type Response struct {
	Text     string `json:"text"`
	Segments []struct {
		Id               int     `json:"id"`
		Seek             int     `json:"seek"`
		Start            float64 `json:"start"`
		End              float64 `json:"end"`
		Text             string  `json:"text"`
		Tokens           []int   `json:"tokens"`
		Temperature      float64 `json:"temperature"`
		AvgLogprob       float64 `json:"avg_logprob"`
		CompressionRatio float64 `json:"compression_ratio"`
		NoSpeechProb     float64 `json:"no_speech_prob"`
	} `json:"segments"`
	Language string `json:"language"`
}

type Config struct {
	// PhraseList      []string `json:"phrase_list"`
	MaxAlternatives int `json:"max_alternatives"`
	SampleRate      int `json:"sample_rate"`
}

type Body struct {
	Config Config `json:"config"`
}

type RemoteRecognizer struct {
	response    chan Response
	status      chan string
	remote      url.URL
	writeBuffer chan bytes.Buffer
	done        chan bool
	state       chan bool
}

type Recognizer interface {
	Connect(host string) (chan bool, error)
	Listen() error
}

func NewRecognizer(response chan Response, status chan string, speaking *bool) Recognizer {
	rr := RemoteRecognizer{
		response:    response,
		status:      status,
		writeBuffer: make(chan bytes.Buffer, 3),
		state:       make(chan bool),
	}

	return &rr

}

func (r *RemoteRecognizer) sendChunk(in []int16) (err error) {
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.LittleEndian, in)
	if err != nil {
		return err
	}
	t := time.NewTimer(time.Millisecond * 500)
	select {
	case r.writeBuffer <- buf: // Put 2 in the channel unless it is full
		t.Stop()
		return nil
	case <-t.C:
		return fmt.Errorf("timed out sending chunk")
	}

}

func (r *RemoteRecognizer) close() error {
	t := time.NewTimer(time.Millisecond * 500)
	select {
	case r.writeBuffer <- bytes.Buffer{}: // Put 2 in the channel unless it is full
		t.Stop()
		return nil
	case <-t.C:
		return fmt.Errorf("timed out closing")
	}
}

func (r *RemoteRecognizer) listen(remote url.URL, ready chan bool) (err error) {
	var conn *websocket.Conn
	conn, _, err = websocket.DefaultDialer.Dial(remote.String(), nil)
	if err != nil {
		return err
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		r.status <- "closed"
		r.done <- true
		log.Event("Disconnecting!")
		return nil
	})

	ready <- true
	for {
		select {
		case buffer := <-r.writeBuffer:
			if buffer.Len() == 0 {
				log.Event("Closing!")
				r.status <- "processing"
				err = conn.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
				if err != nil {
					r.status <- "disconnected"
					return err
				}
				resp := Response{}
				err = conn.ReadJSON(&resp)
				if err != nil {
					return err
				}
				r.status <- "idle"
				r.response <- resp
				return nil
			} else {
				log.Event("Writing!")
				r.status <- "listening"
				err = conn.WriteMessage(websocket.BinaryMessage, buffer.Bytes())
				if err != nil {
					return err
				}
			}
		case <-r.done:
			return nil
		}
	}
}

func (r *RemoteRecognizer) Connect(host string) (doneChan chan bool, err error) {
	r.remote = url.URL{Scheme: "ws", Host: host + ":" + "2700", Path: ""}

	r.done = make(chan bool)

	return r.done, nil
}

const (
	SampleRate = 16000
	BufferSize = 8196
)

func (r *RemoteRecognizer) Listen() error {
	err := portaudio.Initialize()
	if err != nil {
		return err
	}

	buffer := make([]int16, BufferSize)
	stream, err := portaudio.OpenDefaultStream(1, 0, SampleRate, BufferSize, buffer)
	if err != nil {
		return err
	}

	var (
		listening  bool
		quiet      bool
		last       float64
		quietBegin time.Time
	)

	listening = false

	err = stream.Start()

	if err != nil {
		return err
	}

	ready := make(chan bool)

	go func() {
		for {
			log.Event("Connecting!")
			err = r.listen(r.remote, ready)
			if err != nil {
				log.Err(err)
			}
			time.Sleep(time.Second * 5)
		}
	}()

	res := <-ready
	if !res {
		return fmt.Errorf("too many open streams")
	}

	defer func() {
		err = stream.Stop()
		if err != nil {
			log.Err(err)
			return
		}
		err = stream.Close()
		if err != nil {
			return
		}
		err = portaudio.Terminate()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	d := NewDetector(BufferSize)

	for {
		err = stream.Read()
		if err != nil {
			return err
		}

		delta := d.Detect(buffer)

		if last == 0 {
			last = delta
			continue
		}

		if listening {

			if delta*1.75 <= last {
				if !quiet {
					quietBegin = time.Now()
				} else if time.Since(quietBegin) >= time.Second {
					err = r.close()
					if err != nil {
						return err
					}
				}
				quiet = true
			} else {
				err = r.sendChunk(buffer)
				if err != nil {
					log.Err(err)
				}
				quiet = false
				last = delta
			}

		} else {

			if delta >= last*1.75 {

				listening = true
				quiet = false

				err = r.sendChunk(buffer)
				if err != nil {
					log.Err(err)
				}

			}
			last = delta
		}

	}
}
