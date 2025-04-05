package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"romansin312.wt-web/internal/data"
	roomssyncer "romansin312.wt-web/internal/rooms_syncer"
	"romansin312.wt-web/internal/workers"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config     config
	logger     *log.Logger
	models     data.Models
	roomSyncer roomssyncer.RoomSyncer
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Api server port")
	flag.StringVar(&cfg.env, "env", "developemnt", "Environment")

	var dbName string
	flag.StringVar(&dbName, "dbName", "", "Database Name")

	var pgLogin string
	flag.StringVar(&pgLogin, "pgLogin", "", "Database Login")

	var pgPassword string
	flag.StringVar(&pgPassword, "pgPassword", "", "Database Password")

	flag.Parse()

	if dbName == "" {
		panic("Database Name is not provided")
	}
	if dbName == "" {
		panic("Database Login is not provided")
	}
	if dbName == "" {
		panic("Database Password is not provided")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", pgLogin, pgPassword, dbName)

	db, err := initDB(connStr)
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config:     cfg,
		logger:     logger,
		models:     data.NewModels(db),
		roomSyncer: roomssyncer.CreateSyncer(),
	}

	go workers.StartConnectionsKicker(&app.roomSyncer)
	go workers.StartRoomsKicker(&app.models)
	go workers.StartRoomsSyncerWorker(&app.roomSyncer)

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func initDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
