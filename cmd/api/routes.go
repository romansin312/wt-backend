package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/room/:roomId/action", app.actionHandler)
	router.HandlerFunc(http.MethodGet, "/room/:roomId/subscribe", app.subscribeHandler)

	return router
}
