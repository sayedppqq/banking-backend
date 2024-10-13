package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/token"
	"github.com/sayedppqq/banking-backend/util"
	"net/http"
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

func (server *Server) helloServer(ctx *gin.Context) {
	fmt.Println("hello...")
	ctx.JSON(http.StatusOK, "sayed")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
