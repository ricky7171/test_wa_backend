package middleware

import (
	"fmt"
	"net/http"
	"strings"

	helper "wa/helpers"

	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. ambil token dari requestnya user. Diambil dari headernya
		plainToken := c.Request.Header.Get("Authorization")

		//2. cek, kalau tokennya kosong, berarti return "no authorization header provided"
		if plainToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		//3. get tokennya saja (tanpa Bearer )
		splitToken := strings.Split(plainToken, "Bearer ")
		reqToken := splitToken[1]

		//4. ubah token jadi tipe signedDetails
		claims, err := helper.ValidateToken(reqToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		//set value name, phone,uid di context
		c.Set("name", claims.Name)
		c.Set("phone", claims.Phone)
		c.Set("user_id", claims.Uid)
		c.Next()

	}
}
