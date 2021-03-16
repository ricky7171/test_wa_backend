package routes

import (
	"github.com/gin-gonic/gin"
)

func ViewRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	incomingRoutes.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	incomingRoutes.GET("/home", func(c *gin.Context) {
		c.HTML(200, "home.html", nil)
	})

	incomingRoutes.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
}
