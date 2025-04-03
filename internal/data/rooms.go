package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	Id          uuid.UUID
	VideoUrl    string
	CreatedDate time.Time
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

func (r RoomModel) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM "rooms"
		WHERE "Id" = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return RoomNotFoundError
	}
	return nil

}

func (r RoomModel) Get(id uuid.UUID) (*Room, error) {
	query := `
		select "Id", "VideoUrl", "CreatedDate"
		from "rooms"
		where "Id" = $1
	`

	var room Room

	err := r.DB.QueryRow(query, id).Scan(&room.Id, &room.VideoUrl, &room.CreatedDate)

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

func (r RoomModel) GetOlderThan(endDate time.Time) ([]*Room, error) {
	query := `
		select "Id", "VideoUrl", "CreatedDate"
		from "rooms"
		where "CreatedDate" < $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, endDate)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rooms := []*Room{}

	for rows.Next() {
		var room Room
		err := rows.Scan(
			&room.Id,
			&room.VideoUrl,
			&room.CreatedDate,
		)

		if err != nil {
			return nil, err
		}

		rooms = append(rooms, &room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
