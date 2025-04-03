package roomssyncer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomSyncer struct {
	ClientsToRoom map[*websocket.Conn]uuid.UUID
	Upgrader      websocket.Upgrader
}

func CreateSyncer() RoomSyncer {
	return RoomSyncer{
		ClientsToRoom: make(map[*websocket.Conn]uuid.UUID),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

type ActionMessage struct {
	SenderUserId int32
	Timestamp    int64
	ActionType   int32
	ActionInfo   string
}

func (syncer *RoomSyncer) SyncRoom(roomId uuid.UUID, message *ActionMessage) {

	for conn := range syncer.ClientsToRoom {
		if syncer.ClientsToRoom[conn] == roomId {

			sendingMessage, err := json.Marshal(message)
			if err == nil {
				conn.WriteMessage(websocket.TextMessage, sendingMessage)
			}
		}
	}
}

func (syncer *RoomSyncer) AddConnection(roomId uuid.UUID, w http.ResponseWriter, r *http.Request) {
	conn, _ := syncer.Upgrader.Upgrade(w, r, nil)
	syncer.ClientsToRoom[conn] = roomId
}

func (syncer *RoomSyncer) StartConnectionsKicker() {
	var pingRoomAttempts map[uuid.UUID]int32 = make(map[uuid.UUID]int32)
	for {
		var connectionsToRemove []*websocket.Conn
		for conn := range syncer.ClientsToRoom {
			println("Ping connection")
			roomId := syncer.ClientsToRoom[conn]
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

			if pingRoomAttempts[syncer.ClientsToRoom[conn]] > 3 {
				println("Attempts is more than 3, closing connection")
				conn.Close()
				connectionsToRemove = append(connectionsToRemove, conn)
				delete(pingRoomAttempts, roomId)
			}
		}

		fmt.Printf("Number of connections to remove is %d\n", len(connectionsToRemove))
		for _, conn := range connectionsToRemove {
			delete(syncer.ClientsToRoom, conn)
		}

		time.Sleep(5 * time.Second)
	}
}
