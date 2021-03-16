package routes

import (
	"wa/hub"

	"github.com/gin-gonic/gin"
)

func WebSocketRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/ws/:user_id", func(c *gin.Context) {
		userID := c.Param("user_id")
		hub.ServeWs(c.Writer, c.Request, userID)
	})
}
