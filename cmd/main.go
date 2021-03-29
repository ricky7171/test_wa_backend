package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	db "github.com/ricky7171/test_wa_backend/internal/database"
	"github.com/ricky7171/test_wa_backend/internal/hub"
	"github.com/ricky7171/test_wa_backend/internal/middleware"
	routes "github.com/ricky7171/test_wa_backend/internal/routes"
	"go.mongodb.org/mongo-driver/mongo"

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

	//init .env
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//init database
	var dbInstance *mongo.Database = db.DBinstance()

	//run hub to listen data chat websocket on channel
	go hub.MainHub.Run(dbInstance)

	//make router instance
	router := gin.New()

	//load all html
	router.LoadHTMLGlob("web/*")

	// Middleware that used to log all request on terminal
	//router.Use(gin.Logger())

	// Middleware that used to setting CORS
	router.Use(middleware.CORSMiddleware())

	//front-end router
	routes.ViewRoutes(router)

	//authorized router
	authorized := router.Group("/")
	authorized.Use(middleware.Authentication())
	{
		//- chat router
		routes.ChatRoutes(authorized, dbInstance)
	}

	//unauthorized router
	//- auth router
	routes.AuthRoutes(router, dbInstance)

	//run server
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
