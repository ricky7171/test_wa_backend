package routes

import (
	controller "wa/controllers"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.GET("/ws/:userId", controller.ConnectWs())
	incomingRoutes.GET("/api/chat/:contactId/:lastId", controller.GetChat())
	incomingRoutes.POST("/api/new_chat", controller.NewMessage())
	incomingRoutes.GET("/api/contact", controller.GetContact())
}
