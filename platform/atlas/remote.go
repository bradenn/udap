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
	response chan Response
	status   chan string

	writeBuffer    chan bytes.Buffer
	closeBuffer    chan bool
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

func NewRecognizer(response chan Response, status chan string) Recognizer {
	return &RemoteRecognizer{
		response:    response,
		status:      status,
		writeBuffer: make(chan bytes.Buffer),
		closeBuffer: make(chan bool),
		last:        0.0,
		listening:   false,
		quiet:       true,
		threshold:   1.75,
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
	default:
		return nil
	}

	return nil

}

func (r *RemoteRecognizer) close() {

	select {
	case r.closeBuffer <- true: // Put 2 in the channel unless it is full
	default:
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
		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			return
		}
		conn.Close()
	}()

	conn.SetCloseHandler(func(code int, text string) error {
		r.status <- "closed"
		return nil
	})

	for {
		select {
		case buffer := <-r.writeBuffer:
			r.status <- "listening"
			start := time.Now()
			err = conn.WriteMessage(websocket.BinaryMessage, buffer.Bytes())
			if err != nil {
				return err
			}
			_, data, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			fmt.Printf("Sent chunk, took: (%s)\n", time.Since(start))
		case <-r.closeBuffer:
			r.status <- "processing"
			err = conn.WriteMessage(websocket.BinaryMessage, []byte("{\"eof\" : 1}"))
			if err != nil {
				r.status <- "disconnected"
				return err
			}
			var text []byte
			_, text, err = conn.ReadMessage()
			if err != nil {
				r.status <- "disconnected"
				return err
			}

			resp := Response{}
			err = json.Unmarshal(text, &resp)
			if err != nil {
				return err
			}
			r.response <- resp
			r.status <- "idle"
			return nil
		case <-done:
			fmt.Println("Exiting atlas remote recognizer")
			return nil
		}
	}
}

func (r *RemoteRecognizer) Connect(host string) (doneChan chan bool, err error) {
	remote := url.URL{Scheme: "ws", Host: host + ":" + "2700", Path: ""}
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Exiting atlas recognizer")
				return
			default:
				err = r.listen(remote)
				if err != nil {
					log.Err(err)
					return
				}
			}

		}
	}()

	return done, nil
}

func (r *RemoteRecognizer) Listen() error {
	err := portaudio.Initialize()
	if err != nil {
		fmt.Println(err)
	}

	defer func() {
		err = portaudio.Terminate()
		if err != nil {
			fmt.Println(err)
		}
	}()

	in := make([]int16, 8196)

	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	if err != nil {
		return err
	}

	var (
		listening  bool
		quiet      bool
		last       float64
		quietBegin time.Time
	)
	fmt.Println("Beginning to listen:")
	err = stream.Start()
	if err != nil {
		return err
	}
	d := NewDetector(len(in))

	for {
		err = stream.Read()
		if err != nil {
			return err
		}

		delta := d.Detect(in)

		if r.last == 0 {
			r.last = delta
			continue
		}

		if listening {
			fmt.Printf("LEVEL: %.4f\n", delta)
			err = r.sendChunk(in)
			if err != nil {
				return err
			}
			if delta*r.threshold <= last {
				if !quiet {
					quietBegin = time.Now()
				} else {
					diff := time.Since(quietBegin)
					if diff >= time.Millisecond*250 {
						// r.listening = false
						r.close()
						listening = false
						continue
					}
				}
				quiet = true

			} else {
				quiet = false
				last = delta

			}
		} else {
			if delta >= last*r.threshold {
				listening = true
				err = r.sendChunk(in)
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
