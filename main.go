package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sayedppqq/banking-backend/api"
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not create pg connection pool", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("can not setup new server", err)
	}

	err = server.Run(config.HostAddress)
	if err != nil {
		log.Fatal("can not start server", err)
	}
}
