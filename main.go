package main

import (
	"wa/hub"
	routes "wa/routes"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func main() {

	//run hub to listen data chat websocket on channel
	go hub.MainHub.Run()

	//make router instance
	router := gin.New()

	//load all html
	router.LoadHTMLGlob("views/*")

	// Middleware that used to log all request on terminal
	router.Use(gin.Logger())

	//front-end router
	routes.ViewRoutes(router)

	//websocket router
	routes.WebSocketRouter(router)

	//REST API router
	routes.AuthRoutes(router)
	routes.ChatRoutes(router)

	//run server on port 8080
	router.Run("0.0.0.0:8080")
}
