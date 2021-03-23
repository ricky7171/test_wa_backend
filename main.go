package main

import (
	"os"
	"wa/hub"
	"wa/middleware"
	routes "wa/routes"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.Default())

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

	//run server
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
