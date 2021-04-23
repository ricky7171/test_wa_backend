package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ricky7171/test_wa_backend/internal/websocket"
)

type Socket struct {
}

func NewSocketController() *Socket {
	return &Socket{}
}

func (s *Socket) ConnectWs() gin.HandlerFunc {

	return func(c *gin.Context) {
		userID := c.Param("userId")
		websocket.ServeWs(c.Writer, c.Request, userID)
	}
}
