package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"romansin312.wt-web/internal/data"
	roomssyncer "romansin312.wt-web/internal/rooms_syncer"
)

func parseRoomId(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())

	idStr := params.ByName("roomId")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return uuid.Nil, errors.New("roomId is not provided")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (app *application) actionHandler(w http.ResponseWriter, r *http.Request) {
	message := roomssyncer.ActionMessage{}

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go app.roomSyncer.SyncRoom(&message)

	fmt.Printf("Received message: %v\n", r.Body)
}

func (app *application) createHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id       uuid.UUID
		VideoUrl string
	}

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&input)

	room := data.Room{
		Id:       input.Id,
		VideoUrl: input.VideoUrl,
	}
	err := app.models.Rooms.Insert(&room)

	if err != nil {
		fmt.Print(err)
	}

	response, err := room.Id.MarshalText()
	if err != nil {
		fmt.Print(err)
	}
	w.Write(response)
}

func (app *application) getRoomHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseRoomId(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	room, err := app.models.Rooms.Get(id)
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(err, data.RoomNotFoundError):
			http.NotFound(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)

		}

		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(room)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseRoomId(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	app.roomSyncer.AddConnection(id, w, r)

	fmt.Printf("Client has been subscribed on room %s\n", id)
}
