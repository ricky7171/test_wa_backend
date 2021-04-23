package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ricky7171/test_wa_backend/internal/controller"
	db "github.com/ricky7171/test_wa_backend/internal/database"
	"github.com/ricky7171/test_wa_backend/internal/helper"
	"github.com/ricky7171/test_wa_backend/internal/middleware"
	routes "github.com/ricky7171/test_wa_backend/internal/routes"
	"github.com/ricky7171/test_wa_backend/internal/usecase/chat"
	"github.com/ricky7171/test_wa_backend/internal/usecase/contact"
	"github.com/ricky7171/test_wa_backend/internal/usecase/user"
	"github.com/ricky7171/test_wa_backend/internal/websocket"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func main() {

	//1. init .env
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//2. init database
	var dbInstance *mongo.Database = db.DBinstance()
	defer dbInstance.Client().Disconnect(context.TODO())

	//3. init repository pattern && helper
	userRepo := user.NewMongoRepository(dbInstance)
	contactRepo := contact.NewMongoRepository(dbInstance)
	chatRepo := chat.NewMongoRepository(dbInstance)
	tokenHelper := helper.NewTokenJwt()

	//4. init service pattern
	userService := user.NewService(userRepo, tokenHelper)
	contactService := contact.NewService(contactRepo)
	chatService := chat.NewService(chatRepo, userRepo, contactRepo)

	//5. init middleware
	authMiddleware := middleware.NewAuthController(*userService)
	corsMiddleware := middleware.NewCorsController()

	//6. init controller
	userController := controller.NewUserController(*userService)
	contactController := controller.NewContactController(*contactService)
	chatController := controller.NewChatController(*chatService)
	socketController := controller.NewSocketController()

	//7. init router instance and laod html files
	router := gin.New()
	router.LoadHTMLGlob("web/*")

	//8. middleware that used to log all request on terminal
	//router.Use(gin.Logger())

	//9. setup Middleware that used to setting CORS
	router.Use(corsMiddleware.Check())

	//10. setup front-end router
	routes.ViewRoutes(router)

	//11. setup authorized router that contains : chat router
	authorized := router.Group("/")
	authorized.Use(authMiddleware.Check())
	{
		routes.ChatRoutes(authorized, chatController, contactController, socketController)
	}

	//12. setup unauthorized router that contains : auth router
	routes.AuthRoutes(router, userController)

	//13. run hub to listen data chat websocket on channel
	go websocket.MainHub.Run(chatService)

	//14. run server
	var port string
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)

}
