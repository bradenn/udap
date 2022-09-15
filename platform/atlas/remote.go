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
	"sync"
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
	response    chan Response
	status      chan string
	speaking    *bool
	remote      url.URL
	writeBuffer chan bytes.Buffer
	closeBuffer chan bool
	done        chan bool
	listening   bool
	quiet       bool
	threshold   float64
	isListening bool
	mutex       sync.RWMutex
	state       chan bool
}

func (r *RemoteRecognizer) requestState(isListening bool) error {
	select {
	case r.state <- isListening:
		return nil
	case <-time.After(time.Millisecond * 100):
		return fmt.Errorf("timedout")
	}
}

func (r *RemoteRecognizer) checkListening() bool {
	value := false
	r.mutex.RLock()
	value = r.isListening
	r.mutex.RUnlock()
	return value
}

func (r *RemoteRecognizer) listenState() {
	for {
		select {
		case res := <-r.state:
			r.mutex.Lock()
			r.isListening = res
			r.mutex.Unlock()
		}
	}
}

type Recognizer interface {
	Connect(host string) (chan bool, error)
	Listen() error
}

func NewRecognizer(response chan Response, status chan string, speaking *bool) Recognizer {
	rr := RemoteRecognizer{
		response:    response,
		status:      status,
		writeBuffer: make(chan bytes.Buffer, 8),
		closeBuffer: make(chan bool),
		state:       make(chan bool),
		mutex:       sync.RWMutex{},
		isListening: false,
		listening:   false,
		quiet:       true,
		threshold:   1.75,
		speaking:    speaking,
	}

	go rr.listenState()

	return &rr

}

func (r *RemoteRecognizer) sendChunk(in []int16) (err error) {
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.LittleEndian, in)
	if err != nil {
		return err
	}

	select {
	case r.writeBuffer <- buf: // Put 2 in the channel unless it is full
		return nil
	case <-time.After(time.Millisecond * 400):
		return fmt.Errorf("timed out sending chunk")
	}

}

func (r *RemoteRecognizer) close() error {
	if !r.checkListening() {
		return fmt.Errorf("close send failed, not listening")
	}
	select {
	case r.closeBuffer <- true: // Put 2 in the channel unless it is full
		return nil
	case <-time.After(time.Millisecond * 400):
		return fmt.Errorf("timed out closing")
	}
}

func (r *RemoteRecognizer) listen(remote url.URL, ready chan bool) (err error) {
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

	marshal, err := json.Marshal(e)
	if err != nil {
		ready <- false
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		ready <- false
		return err
	}

	defer func() {
		err = r.requestState(false)
		if err != nil {
			log.Err(fmt.Errorf("failed to mark listening state"))
			return
		}
		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure,
			"done"))
		if err != nil {
			log.Err(err)
			return
		}
		err = conn.Close()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		r.status <- "closed"
		r.done <- true
		return nil
	})
	err = r.requestState(true)
	if err != nil {
		ready <- false
		return err
	}
	ready <- true
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
				return err
			}
			r.status <- "idle"
			r.response <- resp
			return nil
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
	BufferSize = 4000
)

func (r *RemoteRecognizer) Listen() error {
	err := portaudio.Initialize()
	if err != nil {
		return err
	}

	buffer := make([]int16, 4000)
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
	err = stream.Start()

	if err != nil {
		return err
	}

	defer func() {
		err = stream.Stop()
		if err != nil {
			log.Err(err)
			return
		}
		err = portaudio.Terminate()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	timeout := time.Now()

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

			if delta*2.0 <= last {
				if !quiet {
					quietBegin = time.Now()
				} else if time.Since(quietBegin) >= time.Millisecond*250 {
					err = r.close()
					if err != nil {
						return err
					}
					return nil
				} else if time.Since(timeout) >= time.Second*8 {
					err = r.close()
					if err != nil {
						return err
					}
					return nil
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

			if delta >= last*2.25 {
				ready := make(chan bool)

				go func() {
					err = r.listen(r.remote, ready)
					if err != nil {
						log.Err(err)
					}
				}()

				res := <-ready
				if !res {
					return fmt.Errorf("too many open streams")
				}

				listening = true
				quiet = false

				timeout = time.Now()

				err = r.sendChunk(buffer)
				if err != nil {
					log.Err(err)
				}

			}
			last = delta
		}

	}
}
