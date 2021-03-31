package routes

import (
	controller "github.com/ricky7171/test_wa_backend/internal/controllers"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine, dbInstance *mongo.Database) {
	incomingRoutes.POST("/api/auth/register", controller.Register(dbInstance))
	incomingRoutes.POST("/api/auth/login", controller.Login(dbInstance))
	incomingRoutes.POST("/api/auth/refresh-token", controller.RefreshToken(dbInstance))
	incomingRoutes.POST("/api/auth/check-token", controller.CheckToken(dbInstance))
}
