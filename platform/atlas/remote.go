// Copyright (c) 2022 Braden Nicholson

package atlas

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
	"udap/internal/log"
)

type Response struct {
	Result []struct {
		Conf  float64 `json:"conf"`
		End   float64 `json:"end"`
		Start float64 `json:"start"`
		Word  string  `json:"word"`
	} `json:"result"`
	Text string `json:"text"`
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
	response       chan Response
	status         chan string
	speaking       *bool
	remote         url.URL
	writeBuffer    chan bytes.Buffer
	closeBuffer    chan bool
	done           chan bool
	listeningSince time.Time
	quietTime      time.Time
	last           float64
	listening      bool
	quiet          bool
	threshold      float64
}

type Recognizer interface {
	Connect(host string) (chan bool, error)
	Listen() error
}

func NewRecognizer(response chan Response, status chan string, speaking *bool) Recognizer {
	return &RemoteRecognizer{
		response:    response,
		status:      status,
		writeBuffer: make(chan bytes.Buffer),
		closeBuffer: make(chan bool),
		last:        0.0,
		listening:   false,
		quiet:       true,
		threshold:   1.75,
		speaking:    speaking,
	}
}

func (r *RemoteRecognizer) sendChunk(in []int16) (err error) {
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.LittleEndian, in)
	if err != nil {
		return err
	}

	select {
	case r.writeBuffer <- buf: // Put 2 in the channel unless it is full
	case <-time.After(time.Millisecond * 250):
		fmt.Println("Timed out sending!")
		return nil
	}

	return nil

}

func (r *RemoteRecognizer) close() {

	select {
	case r.closeBuffer <- true: // Put 2 in the channel unless it is full
	case <-time.After(time.Millisecond * 250):
		fmt.Println("Timed out closing!")
		return
	}

}

func (r *RemoteRecognizer) listen(remote url.URL) (err error) {

	var conn *websocket.Conn
	conn, _, err = websocket.DefaultDialer.Dial(remote.String(), nil)
	if err != nil {
		return err
	}
	e := Body{
		Config: Config{
			SampleRate: 16000,
		},
	}

	done := make(chan bool)
	marshal, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		return err
	}

	defer func() {
		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure,
			"done"))
		if err != nil {
			return
		}
		conn.Close()
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		r.status <- "closed"
		done <- true
		return nil
	})

	for {
		select {
		case buffer := <-r.writeBuffer:
			r.status <- "listening"
			err = conn.WriteMessage(websocket.BinaryMessage, buffer.Bytes())
			if err != nil {
				return err
			}
			_, _, err = conn.ReadMessage()
			if err != nil {
				return err
			}
		case <-r.closeBuffer:
			r.status <- "processing"
			err = conn.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
			if err != nil {
				r.status <- "disconnected"
				return err
			}
			resp := Response{}
			err = conn.ReadJSON(&resp)
			if err != nil {
				return
			}
			r.response <- resp
			r.status <- "idle"
			return nil
		case <-r.done:
			fmt.Println("Exiting atlas remote recognizer")
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
	// VadMode vad mode
	VadMode = 0
	// SampleRate sample rate
	SampleRate = 16000
	// BitDepth bit depth
	BitDepth = 16
	// FrameDuration frame duration
	FrameDuration = 20
)

func (r *RemoteRecognizer) Listen() error {
	err := portaudio.Initialize()
	if err != nil {
		return err
	}

	defer func() {
		err = portaudio.Terminate()
		if err != nil {
			fmt.Println(err)
		}
	}()

	buffer := make([]int16, 4096)
	stream, err := portaudio.OpenDefaultStream(1, 0, SampleRate, len(buffer), buffer)
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
	fmt.Println("Beginning to listen:")
	listenStart := time.Now()
	err = stream.Start()
	if err != nil {
		return err
	}
	d := NewDetector(len(buffer))
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
			err = r.sendChunk(buffer)
			if err != nil {
				return err
			}
			if delta*2 <= last || time.Since(listenStart) >= time.Second*10 {
				if !quiet {
					quietBegin = time.Now()
				} else {
					diff := time.Since(quietBegin)
					if diff > time.Millisecond*500 {
						// r.listening = false
						r.close()
						break
					}
				}
				quiet = true
			} else {
				quiet = false
				last = delta
			}
		} else {

			if delta >= last*2.5 {
				go func() {
					err = r.listen(r.remote)
					if err != nil {
						log.Err(err)
						return
					}
				}()
				time.Sleep(time.Millisecond * 50)
				listenStart = time.Now()
				quiet = false
				listening = true
				err = r.sendChunk(buffer)
				if err != nil {
					return err
				}

			}
			last = delta
		}

	}

	err = stream.Stop()
	if err != nil {
		return err
	}

	return nil

}
