package workers

import (
	"fmt"
	"time"

	"romansin312.wt-web/internal/data"
)

func StartRoomsKicker(models *data.Models) {
	for {
		roomsToKick, err := models.Rooms.GetOlderThan(time.Now().Add(-time.Hour * 24).UTC())
		if err != nil {
			fmt.Printf("Error on rooms fetching: %v\n", err)
		}

		if len(roomsToKick) == 0 {
			println("No rooms to remove")
		} else {
			for _, room := range roomsToKick {
				fmt.Printf("Removing the room with ID=%s\n", room.Id)
				err := models.Rooms.Delete(room.Id)
				if err != nil {
					fmt.Printf("Error while removing the room: %v", err)
				}
			}
		}

		time.Sleep(time.Hour)
	}
}
