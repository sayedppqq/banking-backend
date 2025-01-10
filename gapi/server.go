package gapi

import (
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/pb"
	"github.com/sayedppqq/banking-backend/token"
	"github.com/sayedppqq/banking-backend/util"
)

type Server struct {
	pb.UnimplementedBankingBackendServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}
