package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var clientsToRoom = make(map[*websocket.Conn]string)

func (app *application) actionHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id := params.ByName("roomId")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	var message struct {
		SenderUserId int32
		Timestamp    int64
		ActionType   int32
		ActionInfo   string
	}

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	for conn := range clientsToRoom {
		if clientsToRoom[conn] == id {

			sendingMessage, err := json.Marshal(message)
			if err == nil {
				conn.WriteMessage(websocket.TextMessage, sendingMessage)
			}
		}
	}

	fmt.Printf("Received message: %v", r.Body)
}

func (app *application) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("roomId")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	clientsToRoom[conn] = id
	fmt.Printf("Client has been subscribed on room %s\n", id)
}
