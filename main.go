package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func main() {
	//setup
	go h.run()
	router := gin.New()
	router.LoadHTMLGlob("views/*")

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.String(http.StatusInternalServerError, fmt.Sprintf("Panic !"))
	}))

	//front-end router
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	router.GET("/home", func(c *gin.Context) {
		c.HTML(200, "home.html", nil)
	})

	router.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	//websocket router
	//sebenernya bukan nge get sih. tapi "pada saat ada client yang konek ke 'topik' roomId ini, maka jadikan client ini sebagai subscriber"
	//jadi kayak mendaftarkan client ke topik nya. Setiap ada update dari topic, maka client yang berlangganan akan dikasi tau.
	//client juga bisa kirim data (dengan catatan sudah terdaftar di list subscriber) ke topic roomId tsb

	router.GET("/ws/:token", func(c *gin.Context) {
		//token := c.Param("token")
		//userID = authenticate(token)
		//serveWs(c.Writer, c.Request, userID)
	})

	//REST API router
	router.POST("/api/auth/login", func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"error":    false,
			"message":  "Success Login",
			"phone":    loginRequest.Phone,
			"password": loginRequest.Password,
		})
	})

	router.POST("/api/auth/register", func(c *gin.Context) {

	})

	router.POST("/api/getContacts", func(c *gin.Context) {

	})

	router.Run("0.0.0.0:8080")
}
