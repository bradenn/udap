// Copyright (c) 2021 Braden Nicholson

package bond

import (
	"encoding/json"
	"fmt"
	"time"
)

type Msg struct {
	Target    string `json:"target"`    // Module, Daemon, Agent, etc.
	Operation string `json:"operation"` // state
	Id        string `json:"id"`        // {state: "on"}
	Body      any    `json:"body"`      // {state: "on"}
	Payload   string `json:"payload"`   // {state: "on"}
	resolved  bool
	Resp      chan Resp
}

func (m *Msg) FromJSON(data json.RawMessage) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Msg) ResponseJSON(body []byte, err error) {
	if m.resolved {
		return
	}

	m.resolved = true
	if err != nil {
		m.Resp <- Resp{
			Success: false,
			Error:   err.Error(),
			Body:    body,
		}
	} else {
		m.Resp <- Resp{
			Success: true,
			Error:   "",
			Body:    body,
		}
	}

}

func (m *Msg) Respond(body any, err error) {
	if err != nil {
		m.Resp <- Resp{
			Success: false,
			Error:   err.Error(),
			Body:    body,
		}
	} else {
		m.Resp <- Resp{
			Success: true,
			Error:   "",
			Body:    body,
		}
	}

}

func (m *Msg) Success() {
	m.Resp <- Resp{
		Success: true,
		Error:   "nil",
		Body:    nil,
	}
}

type Resp struct {
	Success bool   `json:"status"`
	Error   string `json:"error"`
	Body    any    `json:"body"`
}

type Bond struct {
	channel chan Msg
}

func NewBond(channel chan Msg) *Bond {
	return &Bond{channel: channel}
}

func (b *Bond) send(target Msg) (any, error) {
	if b.channel == nil {
		return nil, fmt.Errorf("channel not connected")
	}
	c := make(chan Resp)
	target.Resp = c
	b.channel <- target
	select {
	case res, ok := <-c:
		close(c)
		if !ok {
			return nil, fmt.Errorf("somebody made a fucky-wucky")
		}
		if res.Error != "" {
			return res.Body, fmt.Errorf("command '%s.%s' failed: %s", target.Target, target.Operation, res.Error)
		}
		return res.Body, nil
	case <-time.After(2 * time.Second):
		return nil, fmt.Errorf("timeout: %s.%s (%s) => %s", target.Target, target.Operation, target.Id, target.Body)
	}
}

func (b *Bond) SendIdJSON(target string, operation string, id string, body string) (any, error) {
	return b.send(Msg{
		Target:    target,
		Operation: operation,
		Id:        id,
		Payload:   body,
	})
}
func (b *Bond) SendId(target string, operation string, id string, body any) (any, error) {
	return b.send(Msg{
		Target:    target,
		Operation: operation,
		Id:        id,
		Body:      body,
	})
}

func (b *Bond) CmdJSON(body []byte) (any, error) {
	msg := Msg{}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return nil, err
	}
	return b.send(msg)
}

func (b *Bond) Msg(msg Msg) (any, error) {
	return b.send(msg)
}

func (b *Bond) Send(target string, operation string, body any) (any, error) {
	return b.send(Msg{
		Target:    target,
		Operation: operation,
		Body:      body,
	})
}
