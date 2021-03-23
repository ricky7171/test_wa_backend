package main

import (
	"os"
	"wa/hub"
	"wa/middleware"
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
	//router.Use(gin.Logger())

	//front-end router
	routes.ViewRoutes(router)

	//authorized router
	authorized := router.Group("/")
	authorized.Use(middleware.Authentication())
	{
		//- chat router
		routes.ChatRoutes(authorized)
	}

	//unauthorized router
	//- auth router
	routes.AuthRoutes(router)

	//run server on port 8080
	router.Run(":" + os.Getenv("PORT"))
}
