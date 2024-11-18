package api

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/users/access_token_renew", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)
	server.router = router
}
