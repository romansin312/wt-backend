package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type roomSyncer struct {
	clientsToRoom map[*websocket.Conn]uuid.UUID
	upgrader      websocket.Upgrader
}

func CreateSyncer() roomSyncer {
	return roomSyncer{
		clientsToRoom: make(map[*websocket.Conn]uuid.UUID),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

type actionMessage struct {
	SenderUserId int32
	Timestamp    int64
	ActionType   int32
	ActionInfo   string
}

func (syncer *roomSyncer) syncRoom(roomId uuid.UUID, message *actionMessage) {

	for conn := range syncer.clientsToRoom {
		if syncer.clientsToRoom[conn] == roomId {

			sendingMessage, err := json.Marshal(message)
			if err == nil {
				conn.WriteMessage(websocket.TextMessage, sendingMessage)
			}
		}
	}
}

func (syncer *roomSyncer) addConnection(roomId uuid.UUID, w http.ResponseWriter, r *http.Request) {
	conn, _ := syncer.upgrader.Upgrade(w, r, nil)
	syncer.clientsToRoom[conn] = roomId
}

func (syncer *roomSyncer) startConnectionsKicker() {
	var pingRoomAttempts map[uuid.UUID]int32 = make(map[uuid.UUID]int32)
	for {
		var connectionsToRemove []*websocket.Conn
		for conn := range syncer.clientsToRoom {
			println("Ping connection")
			roomId := syncer.clientsToRoom[conn]
			if pingRoomAttempts[roomId] == 0 {
				pingRoomAttempts[roomId] = 1
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
			if err != nil {
				fmt.Printf("Error occured, an attempt is %d\n", pingRoomAttempts[roomId])
				pingRoomAttempts[roomId]++
			} else {
				pingRoomAttempts[roomId] = 1
			}

			if pingRoomAttempts[syncer.clientsToRoom[conn]] > 3 {
				println("Attempts is more than 3, closing connection")
				conn.Close()
				connectionsToRemove = append(connectionsToRemove, conn)
				delete(pingRoomAttempts, roomId)
			}
		}

		fmt.Printf("Number of connections to remove is %d\n", len(connectionsToRemove))
		for _, conn := range connectionsToRemove {
			delete(syncer.clientsToRoom, conn)
		}

		time.Sleep(5 * time.Second)
	}
}
