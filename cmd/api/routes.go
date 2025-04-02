package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/room/:roomId/action", app.actionHandler)
	router.HandlerFunc(http.MethodGet, "/room/:roomId/subscribe", app.subscribeHandler)
	router.HandlerFunc(http.MethodGet, "/room/:roomId/get", app.getRoomHandler)
	router.HandlerFunc(http.MethodPost, "/rooms/create", app.createHandler)

	return app.enableCORS(router)
}
