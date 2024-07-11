package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/token"
	"github.com/sayedppqq/banking-backend/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
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
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("validCurrency", validCurrency)
		if err != nil {
			return nil, err
		}
	}
	server.setupRouter()
	return server, nil
}
func (server *Server) Run(address string) error {
	return server.router.Run(address)
}
