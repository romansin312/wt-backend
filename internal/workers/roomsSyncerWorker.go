package workers

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	roomssyncer "romansin312.wt-web/internal/rooms_syncer"
)

func StartRoomsSyncerWorker(roomsSyncer *roomssyncer.RoomSyncer) {
	for {
		message := roomsSyncer.GetMessage()
		fmt.Printf("Message is received for roomId=%s\n", message.RoomId)
		roomId := message.RoomId
		connectionsToClient := roomsSyncer.GetConnectionsToClientMap()
		for conn := range connectionsToClient {
			clientInfo := connectionsToClient[conn]
			if clientInfo.RoomId == roomId && clientInfo.UserId != message.SenderUserId {
				sendingMessage, err := json.Marshal(message)
				if err == nil {
					conn.WriteMessage(websocket.TextMessage, sendingMessage)
				}
			}
		}
	}
}
