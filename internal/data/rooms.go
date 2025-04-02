package data

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type Room struct {
	Id       uuid.UUID
	VideoUrl string
}

type RoomModel struct {
	DB *sql.DB
}

var RoomNotFoundError = errors.New("room is not found")

func (r RoomModel) Insert(room *Room) error {
	query := `
		INSERT INTO "rooms" ("Id", "VideoUrl")
		VALUES ($1, $2)
		RETURNING "Id"
	`

	args := []any{uuid.New(), room.VideoUrl}

	return r.DB.QueryRow(query, args...).Scan(&room.Id)
}

func (r RoomModel) Get(id uuid.UUID) (*Room, error) {
	query := `
		select "Id", "VideoUrl"
		from "rooms"
		where "Id" = $1
	`

	var room Room

	err := r.DB.QueryRow(query, id).Scan(&room.Id, &room.VideoUrl)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, RoomNotFoundError
		default:
			return nil, err
		}
	}

	return &room, nil
}
