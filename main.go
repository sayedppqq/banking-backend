package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sayedppqq/banking-backend/api"
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/util"
	"log"
)

func main() {
	fmt.Println("-----------Starting backend server------------------")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	fmt.Println("111111111111111111")

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not create pg connection pool", err)
	}
	fmt.Println("2222222222222222222222222222")

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("can not setup new server", err)
	}
	fmt.Println("333333333333333333333333", config)

	err = server.Run(config.HostAddress)
	fmt.Println("4444444444444444444", err)
	if err != nil {
		fmt.Println("here..............", err)
		log.Fatal("can not start server", err)
	}
	fmt.Println("server started!!!")
}
