package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"udap/config"
	"udap/server"
)

func main() {
	// Load ENV variables from .env, or context
	config.Init()
	RunDemo()
	// Establish server structure, connect to database
	srv, err := server.New()
	if err != nil {
		log.Fatalln(err)
	}
	uuid, err := server.SignUUID("613eb1713bebb1ebe21c6403")
	if err != nil {
		return
	}
	fmt.Println(uuid)
	// Begin routes requiring jwt authentication
	srv.RouteSecure("/endpoints", RouteEndpoints)
	srv.RouteSecure("/instances", RouteInstances)
	srv.RouteSecure("/modules", RouteModules)
	srv.Router().HandleFunc("/echo", handler)
	// Run the server indefinitely
	err = srv.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

var upg = websocket.Upgrader{}

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
