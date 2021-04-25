package controller

import (
	"net/http"

	"github.com/ricky7171/test_wa_backend/internal/failure"
	"github.com/ricky7171/test_wa_backend/internal/helper"

	"github.com/gin-gonic/gin"

	"github.com/ricky7171/test_wa_backend/internal/usecase/user"
)

type User struct {
	userService user.Service
}

func NewUserController(userService user.Service) *User {
	return &User{
		userService: userService,
	}
}

func (u *User) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. make input struct to store request body
		var input struct {
			Name     string `json:"name"`
			Phone    string `json:"phone"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//2. validate input & create user
		newIDUser, err := u.userService.CreateUser(input.Name, input.Phone, input.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//3. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", map[string]string{
			"InsertedID": newIDUser,
		}))
		c.Abort()
		return

	}
}

func (u *User) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. make input struct to store request body
		var input struct {
			Phone    string `json:"phone"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", failure.ErrCannotReadJson().Error()))
			c.Abort()
			return
		}

		//2. validate input
		if input.Phone == "" || input.Password == "" {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", failure.ErrFieldRequired("phone", "password").Error()))
			c.Abort()
			return
		}

		//3. Authenticate with phone & password from request
		userFound, err := u.userService.Authenticate(input.Phone, input.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. remove password value
		userFound.Password = ""

		//5. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", userFound))

	}
}

//used to refresh access token that has been expired
func (u *User) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		//1. read input
		var input struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//2. validate input
		if input.RefreshToken == "" {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", "Refresh token is required"))
			c.Abort()
			return
		}

		//3. refresh token
		user, err := u.userService.RefreshToken(input.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. remove password value
		user.Password = ""

		//5. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", user))

	}
}

//check token valid or not
func (u *User) CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. read input
		var input struct {
			Token string `json:"token"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//2. validate input
		if input.Token == "" {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", "Token is required"))
			c.Abort()
			return
		}

		//3. check token
		user, err := u.userService.CheckToken(input.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. remove password value
		user.Password = ""

		//5. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", user))
	}
}
