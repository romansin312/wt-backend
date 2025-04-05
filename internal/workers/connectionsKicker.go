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
		clientsToRoom := syncer.GetConnectionsToRoomMap()
		fmt.Printf("Number of active connections is %d\n", len(clientsToRoom))

		for conn := range clientsToRoom {
			println("Ping connection")
			roomId := clientsToRoom[conn]
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

			if pingRoomAttempts[clientsToRoom[conn]] > 3 {
				println("Attempts is more than 3, closing connection")
				conn.Close()
				delete(pingRoomAttempts, roomId)
				syncer.RemoveConnection(conn)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
