package roomssyncer

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomSyncer struct {
	clientsToRoom        map[*websocket.Conn]uuid.UUID
	Upgrader             websocket.Upgrader
	actionMessagesBuffer chan *ActionMessage
}

func CreateSyncer() RoomSyncer {
	return RoomSyncer{
		clientsToRoom: make(map[*websocket.Conn]uuid.UUID),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		actionMessagesBuffer: make(chan *ActionMessage, 100),
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
	syncer.actionMessagesBuffer <- message
}

func (syncer *RoomSyncer) AddConnection(roomId uuid.UUID, w http.ResponseWriter, r *http.Request) {
	conn, _ := syncer.Upgrader.Upgrade(w, r, nil)
	syncer.clientsToRoom[conn] = roomId
}

func (syncer *RoomSyncer) GetMessage() *ActionMessage {
	return <-syncer.actionMessagesBuffer
}

func (syncer *RoomSyncer) GetConnectionsToRoomMap() map[*websocket.Conn]uuid.UUID {
	result := make(map[*websocket.Conn]uuid.UUID)
	for i := range syncer.clientsToRoom {
		result[i] = syncer.clientsToRoom[i]
	}

	return result
}

func (syncer *RoomSyncer) RemoveConnection(conn *websocket.Conn) {
	delete(syncer.clientsToRoom, conn)
}
