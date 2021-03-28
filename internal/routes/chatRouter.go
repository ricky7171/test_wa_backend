package routes

import (
	controller "github.com/ricky7171/test_wa_backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.GET("/ws/:userId", controller.ConnectWs())
	incomingRoutes.GET("/api/chat/:contactId/:lastId", controller.GetChat())
	incomingRoutes.POST("/api/new-chat", controller.NewMessage())
	incomingRoutes.GET("/api/contact", controller.GetContact())
}
