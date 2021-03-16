package routes

import (
	controller "wa/controllers"
	"wa/middleware"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/api/chat/:room_id/:last_id", controller.GetChat())
	incomingRoutes.POST("/api/new_chat", controller.NewMessage())
	incomingRoutes.GET("/api/contact", controller.GetContact())
}
