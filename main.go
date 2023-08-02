package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/ndavidson19/quanta-backend/api"
	db "github.com/ndavidson19/quanta-backend/db"
	"github.com/ndavidson19/quanta-backend/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: %w", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: %w", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: %w", err)
	}

}
