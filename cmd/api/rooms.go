package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"romansin312.wt-web/internal/data"
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

func (app *application) createHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id       uuid.UUID
		VideoUrl string
	}

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&input)

	room := data.Room{
		Id:       input.Id,
		VideoUrl: input.VideoUrl,
	}
	err := app.models.Rooms.Insert(&room)

	if err != nil {
		fmt.Print(err)
	}

	response, err := room.Id.MarshalText()
	if err != nil {
		fmt.Print(err)
	}
	w.Write(response)
}

func (app *application) getRoomHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	idStr := params.ByName("roomId")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room, err := app.models.Rooms.Get(id)
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(err, data.RoomNotFoundError):
			http.NotFound(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)

		}

		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(room)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("roomId")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	clientsToRoom[conn] = id
	fmt.Printf("Client has been subscribed on room %s\n", id)

	go pingConnection(conn)
}

func pingConnection(conn *websocket.Conn) {
	defer conn.Close()

	attempts := 1
	for {
		println("Ping connection")
		err := conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
		if err != nil {
			fmt.Printf("Error occured, an attempt is %d\n", attempts)
			attempts++
		} else {
			attempts = 1
		}

		if attempts > 3 {
			println("Attempts is more than 3, closing connection")
			return
		}

		time.Sleep(5 * time.Second)
	}
}
