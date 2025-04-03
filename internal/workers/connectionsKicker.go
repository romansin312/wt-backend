package workers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	roomssyncer "romansin312.wt-web/internal/rooms_syncer"
)

func StartConnectionsKicker(syncer *roomssyncer.RoomSyncer) {
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
