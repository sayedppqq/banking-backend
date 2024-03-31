package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sayedppqq/banking-backend/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	server.setupRouter()
	return server
}
func (server *Server) Run(address string) error {
	return server.router.Run(address)
}
