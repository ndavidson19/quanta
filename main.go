package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ndavidson19/quanta-backend/api"
	db "github.com/ndavidson19/quanta-backend/db"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5433/go_client?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: %w", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: %w", err)
	}

}
