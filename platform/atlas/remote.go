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
	"udap/platform/vad"
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

func sendChunk(c *websocket.Conn, in []int16) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, in)
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.BinaryMessage, buf.Bytes())
	if err != nil {
		return err
	}

	_, _, err = c.ReadMessage()
	if err != nil {
		return err
	}

	return nil
}

var (
	last      = 0.0
	listening = false
	quiet     = true
	threshold = 1.75
)

type Config struct {
	// PhraseList      []string `json:"phrase_list"`
	MaxAlternatives int `json:"max_alternatives"`
	SampleRate      int `json:"sample_rate"`
}

type Body struct {
	Config Config `json:"config"`
}

func BeginListening(response chan Response, status chan string) error {
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

	in := make([]int16, 4000)

	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	if err != nil {
		return err
	}

	defer stream.Close()

	u := url.URL{Scheme: "ws", Host: "localhost" + ":" + "2700", Path: ""}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

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

	err = c.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		return err
	}

	defer func() {
		// Closing websocket connection
		err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			return
		}
	}()

	err = stream.Start()
	if err != nil {
		return err
	}

	queue := make(chan []int16, 10)

	go func() {
		for val := range queue {
			err = listenSingle(c, val, response, status)
			if err != nil {
				continue
			}
		}
	}()

	for {
		var read int
		read, err = stream.AvailableToRead()
		if err != nil {
			return err
		}
		if read < 4000 {
			continue
		}

		err = stream.Read()
		if err != nil {
			return err
		}

		queue <- in
	}

}

func listenSingle(c *websocket.Conn, in []int16, response chan Response, status chan string) error {
	d := vad.NewDetector(len(in))
	delta := d.Detect(in)

	if last == 0 {
		last = delta
		return nil
	}

	if listening {
		err := sendChunk(c, in)
		if err != nil {
			return err
		}
		if delta*threshold <= last {
			if quiet {
				listening = false
				status <- "processing"

				var text []byte
				err = c.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
				if err != nil {
					return err
				}
				_, text, err = c.ReadMessage()
				if err != nil {
					return err
				}

				resp := Response{}
				err = json.Unmarshal(text, &resp)
				if err != nil {
					return err
				}

				response <- resp
				status <- "idle"

			}
			quiet = true

		} else {
			quiet = false

			last = delta
		}
	} else {
		if delta >= last*threshold*1.25 {
			status <- "listening"
			listening = true
			err := sendChunk(c, in)
			if err != nil {
				return err
			}
		}
		last = delta
	}
	return nil
}
