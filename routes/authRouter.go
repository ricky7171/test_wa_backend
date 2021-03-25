package routes

import (
	controller "wa/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/auth/register", controller.Register())
	incomingRoutes.POST("/api/auth/login", controller.Login())
	//incomingRoutes.POST("/api/auth/refresh-token", controller.Login())
}
