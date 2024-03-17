package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const (
	driver = "postgres"
	source = "postgresql://root:root@localhost:5432/bank?sslmode=disable"
)

var testStore Store // store for testing purpose.

func TestMain(m *testing.M) {

	connPool, err := pgxpool.New(context.Background(), source)
	if err != nil {
		log.Fatal("can not connect to db")
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
