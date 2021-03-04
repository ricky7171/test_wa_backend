package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	go h.run()
	// fmt.Println("h.run nih")
	router := gin.New()
	router.LoadHTMLFiles("index.html")

	router.GET("/room/:roomId", func(c *gin.Context) {
		// fmt.Println("ada yang masuk room", c.Param("roomId"))
		c.HTML(200, "index.html", nil)
	})

	router.GET("/roomcoba1", func(c *gin.Context) {
		c.JSON(200, "{error:false}")
	})

	//sebenernya bukan nge get sih. tapi "pada saat ada client yang konek ke 'topik' roomId ini, maka jadikan client ini sebagai subscriber"
	//jadi kayak mendaftarkan client ke topik nya. Setiap ada update dari topic, maka client yang berlangganan akan dikasi tau.
	//client juga bisa kirim data (dengan catatan sudah terdaftar di list subscriber) ke topic roomId tsb

	router.GET("/ws/:roomId", func(c *gin.Context) {
		//fmt.Println("ada yang konek ke ws nya room")
		roomId := c.Param("roomId")
		serveWs(c.Writer, c.Request, roomId)
	})

	router.Run("0.0.0.0:8080")
}
