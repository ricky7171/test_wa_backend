package routes

import (
	"github.com/ricky7171/test_wa_backend/internal/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine, userController *controller.User) {
	incomingRoutes.POST("/api/auth/register", userController.Register())
	incomingRoutes.POST("/api/auth/login", userController.Login())
	incomingRoutes.POST("/api/auth/refresh-token", userController.RefreshToken())
	incomingRoutes.POST("/api/auth/check-token", userController.CheckToken())
}
