package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sayedppqq/banking-backend/util"
	"log"
	"os"
	"testing"
)

var testStore Store // store for testing purpose.

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
