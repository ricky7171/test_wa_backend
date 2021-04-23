package routes

import (
	"github.com/ricky7171/test_wa_backend/internal/controller"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(incomingRoutes *gin.RouterGroup, chatController *controller.Chat, contactController *controller.Contact, socketController *controller.Socket) {
	incomingRoutes.GET("/ws/:userId", socketController.ConnectWs())
	incomingRoutes.GET("/api/chat/:contactId/:lastId", chatController.GetChat())
	incomingRoutes.POST("/api/new-chat", chatController.NewMessage())
	incomingRoutes.GET("/api/contact", contactController.GetContact())
}
