package middleware

import (
	"net/http"
	"strings"

	"github.com/ricky7171/test_wa_backend/internal/helper"
	"github.com/ricky7171/test_wa_backend/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	userService user.Service
}

func NewAuthController(userService user.Service) *Auth {
	return &Auth{
		userService: userService,
	}
}

// validates token and authorizes users
func (a *Auth) Check() gin.HandlerFunc {
	return func(c *gin.Context) {

		var plainToken string
		//1. get token from header Authorization (for restful API)
		//or from header Sec-WebSocket-Protocl (for ws connection)
		if c.Request.Header.Get("Authorization") != "" {
			plainToken = c.Request.Header.Get("Authorization")
		} else if c.Query("access_token") != "" {
			plainToken = "Bearer " + c.Query("access_token")
		} else {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "Not authorized !"))
			c.Abort()
			return
		}

		//2. check if token is empty, then return "no authorization header provided"
		if plainToken == "" {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "Not authorized !"))
			c.Abort()
			return
		}

		//3. get token string without "Bearer "
		splitToken := strings.Split(plainToken, "Bearer ")
		reqToken := splitToken[1]

		//4. check token
		user, err := a.userService.CheckToken(reqToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", err))
			c.Abort()

			return
		}

		//5. set value name, phone, ID in gin context
		c.Set("name", user.Name)
		c.Set("phone", user.Phone)
		c.Set("userId", user.ID)
		c.Next()

	}
}
