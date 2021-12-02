// Copyright (c) 2021 Braden Nicholson

package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader websocket.Upgrader

func init() {
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// func (w *Server) handle(writer http.ResponseWriter, request *http.Request) {
// 	fmt.Println(request.Host)
// 	conn, err := Upgrader.Upgrade(writer, request, nil)
// 	if err != nil {
// 		return
// 	}
// 	log.Info("Endpoint attempting to connect")
// 	defer func() {
// 		if err = conn.Close(); err != nil {
// 			log.Err(err)
// 		}
// 	}()
//
// 	if !request.URL.Query().Has("token") {
// 		err = conn.Close()
// 		if err != nil {
// 			return
// 		}
// 	}
//
// 	token := request.URL.Query().Get("token")
//
// 	_, err = AuthToken(token)
// 	if err != nil {
// 		log.Err(err)
// 		return
// 	}
//
// 	for {
// 		var req WSRequest
//
// 		err = conn.ReadJSON(&req)
// 		if err != nil {
// 			return
// 		}
//
// 		req.sender = conn
//
// 		w.socket()
// 	}
// }
//
// func (w *Server) Handle(pattern string, fn func(request WSRequest)) {
// 	// w.manifest[pattern] = fn
// }

type WSRequest struct {
	Path   string `json:"path"`
	Body   string `json:"body"`
	sender *websocket.Conn
}

// path: token.apple.%s.weeds
// args: sauce

func (r *WSRequest) Params() (res map[string]interface{}) {
	err := json.Unmarshal([]byte(r.Body), &res)
	if err != nil {
		return nil
	}
	return res
}

func (r *WSRequest) send(body interface{}) {
	// err := r.Sender.Conn.WriteJSON(body)
	// if err != nil {
	// 	log.ErrF(err, "an error occurred")
	// }
}

func (r *WSRequest) Resolve(body interface{}) {
	r.send(body)
}

func (r *WSRequest) Reject(body interface{}) {
	r.send(body)
}
