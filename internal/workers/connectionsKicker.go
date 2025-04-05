package workers

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	roomssyncer "romansin312.wt-web/internal/rooms_syncer"
)

func StartConnectionsKicker(syncer *roomssyncer.RoomSyncer) {
	var pingUserAttempts map[int32]int32 = make(map[int32]int32)
	for {
		connectionsToClients := syncer.GetConnectionsToClientMap()
		fmt.Printf("Number of active connections is %d\n", len(connectionsToClients))

		for conn := range connectionsToClients {
			println("Ping connection")
			userId := connectionsToClients[conn].UserId
			if pingUserAttempts[userId] == 0 {
				pingUserAttempts[userId] = 1
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
			if err != nil {
				fmt.Printf("Error occured, an attempt is %d\n", pingUserAttempts[userId])
				pingUserAttempts[userId]++
			} else {
				pingUserAttempts[userId] = 1
			}

			if pingUserAttempts[userId] > 3 {
				println("Attempts is more than 3, closing connection")
				conn.Close()
				delete(pingUserAttempts, userId)
				syncer.RemoveConnection(conn)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
