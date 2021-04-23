package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ricky7171/test_wa_backend/internal/helper"
	"github.com/ricky7171/test_wa_backend/internal/usecase/contact"
)

type Contact struct {
	contactService contact.Service
}

func NewContactController(contactService contact.Service) *Contact {
	return &Contact{
		contactService: contactService,
	}
}

//get contact
//contact is peoples who have interacted before with specific user id
func (co *Contact) GetContact() gin.HandlerFunc {
	return func(c *gin.Context) {

		//1. get param userId and user name
		userId := c.GetString("userId")
		userName := c.GetString("name")

		//2. validate input
		if userId == "" || userName == "" {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", "user id or user name cannot be empty"))
			c.Abort()
			return
		}

		//3. get contact
		contacts, err := co.contactService.GetContactByUser(userId, userName)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", contacts))

	}
}
