package routes

import (
	controller "github.com/ricky7171/test_wa_backend/internal/controllers"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(incomingRoutes *gin.RouterGroup, dbInstance *mongo.Database) {
	incomingRoutes.GET("/ws/:userId", controller.ConnectWs())
	incomingRoutes.GET("/api/chat/:contactId/:lastId", controller.GetChat(dbInstance))
	incomingRoutes.POST("/api/new-chat", controller.NewMessage(dbInstance))
	incomingRoutes.GET("/api/contact", controller.GetContact(dbInstance))
}
