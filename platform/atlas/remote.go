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
	Alternatives []struct {
		Confidence float64 `json:"confidence"`
		Result     []struct {
			End   float64 `json:"end"`
			Start float64 `json:"start"`
			Word  string  `json:"word"`
		} `json:"result"`
		Text string `json:"text"`
	} `json:"alternatives"`
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
	conn     *websocket.Conn
	response chan Response
	status   chan string

	writeBuffer    chan bytes.Buffer
	closeBuffer    chan bool
	listeningSince time.Time
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
		writeBuffer: make(chan bytes.Buffer, 20),
		closeBuffer: make(chan bool),
		last:        0.0,
		listening:   false,
		quiet:       true,
		threshold:   2,
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

func (r *RemoteRecognizer) listen(remote url.URL) (err error) {
	r.conn, _, err = websocket.DefaultDialer.Dial(remote.String(), nil)

	e := Body{
		Config: Config{
			MaxAlternatives: 2,
			SampleRate:      16000,
		},
	}

	marshal, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = r.conn.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		return err
	}

	defer func() {
		err = r.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			return
		}
	}()

	r.conn.SetCloseHandler(func(code int, text string) error {
		r.status <- "closed"
		return nil
	})

	for {
		select {
		case buffer := <-r.writeBuffer:
			r.status <- "listening"
			err = r.conn.WriteMessage(websocket.BinaryMessage, buffer.Bytes())
			if err != nil {
				return err
			}
			_, _, err = r.conn.ReadMessage()
			if err != nil {
				return err
			}
		case _ = <-r.closeBuffer:
			r.status <- "processing"
			err = r.conn.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
			if err != nil {
				r.status <- "disconnected"
				return err
			}
			var text []byte
			_, text, err = r.conn.ReadMessage()
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
		}
	}

	return nil
}

func (r *RemoteRecognizer) Connect(host string) (doneChan chan bool, err error) {
	remote := url.URL{Scheme: "ws", Host: host + ":" + "2700", Path: ""}
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				err = r.listen(remote)
				if err != nil {
					log.Err(err)
					continue
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

	in := make([]int16, 3200)

	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	if err != nil {
		return err
	}

	err = stream.Start()
	if err != nil {
		return err
	}

	for {
		err = stream.Read()
		if err != nil {
			return err
		}
		d := NewDetector(len(in))
		delta := d.Detect(in)

		if r.last == 0 {
			r.last = delta
			continue
		}

		if r.listening {

			err = r.sendChunk(in)
			if err != nil {
				return err
			}

			if time.Since(r.listeningSince).Seconds() > 10 {
				r.listening = false
				r.closeBuffer <- true
				fmt.Println("Timedout atlas listening")
				r.quiet = true
				continue
			}

			if delta*r.threshold*2 <= r.last {
				if r.quiet {
					r.listening = false
					r.closeBuffer <- true
				}
				r.quiet = true
			} else {
				r.quiet = false

				r.last = delta
			}
		} else {
			if delta >= r.last*r.threshold*1.5 {
				r.listening = true
				r.listeningSince = time.Now()
				err = r.sendChunk(in)
				if err != nil {
					return err
				}
			}
			r.last = delta
		}

	}

}
