package roomssyncer

import (
	"encoding/json"
	"net/http"

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
	RoomId       uuid.UUID
}

func (syncer *RoomSyncer) SyncRoom(message *ActionMessage) {
	roomId := message.RoomId
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
