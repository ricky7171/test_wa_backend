package middleware

import (
	"fmt"
	"net/http"
	"strings"

	helper "github.com/ricky7171/test_wa_backend/internal/helpers"

	"github.com/gin-gonic/gin"
)

// Authz validates token and authorizes users
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		var plainToken string
		//1. get token from header Authorization (for restful API)
		//or from header Sec-WebSocket-Protocl (for ws connection)
		if c.Request.Header.Get("Authorization") != "" {
			plainToken = c.Request.Header.Get("Authorization")
		} else if c.Query("access_token") != "" {
			plainToken = "Bearer " + c.Query("access_token")
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not authorized !"})
			c.Abort()
			return
		}

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
		c.Set("userId", claims.ID)
		c.Next()

	}
}
