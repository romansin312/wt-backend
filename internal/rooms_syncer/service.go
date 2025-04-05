package roomssyncer

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"romansin312.wt-web/cmd/models"
)

type ClientModel struct {
	UserId int32
	RoomId uuid.UUID
}

type RoomSyncer struct {
	connectionsToClients map[*websocket.Conn]ClientModel
	Upgrader             websocket.Upgrader
	actionMessagesBuffer chan *models.ActionMessage
}

func CreateSyncer() RoomSyncer {
	return RoomSyncer{
		connectionsToClients: make(map[*websocket.Conn]ClientModel),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		actionMessagesBuffer: make(chan *models.ActionMessage, 100),
	}
}

func (syncer *RoomSyncer) SyncRoom(message *models.ActionMessage) {
	syncer.actionMessagesBuffer <- message
}

func (syncer *RoomSyncer) AddConnection(roomId uuid.UUID, userId int32, w http.ResponseWriter, r *http.Request) {
	conn, _ := syncer.Upgrader.Upgrade(w, r, nil)
	syncer.connectionsToClients[conn] = ClientModel{
		UserId: userId,
		RoomId: roomId,
	}
}

func (syncer *RoomSyncer) GetMessage() *models.ActionMessage {
	return <-syncer.actionMessagesBuffer
}

func (syncer *RoomSyncer) GetConnectionsToClientMap() map[*websocket.Conn]ClientModel {
	result := make(map[*websocket.Conn]ClientModel)
	for i := range syncer.connectionsToClients {
		result[i] = syncer.connectionsToClients[i]
	}

	return result
}

func (syncer *RoomSyncer) RemoveConnection(conn *websocket.Conn) {
	disconnectedMessage := models.ActionMessage{
		SenderUserId: syncer.connectionsToClients[conn].UserId,
		Timestamp:    time.Now().UTC().Unix(),
		ActionType:   models.UserDisconnected,
		ActionInfo:   "",
		RoomId:       syncer.connectionsToClients[conn].RoomId,
	}
	syncer.actionMessagesBuffer <- &disconnectedMessage
	delete(syncer.connectionsToClients, conn)
}
