package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ricky7171/test_wa_backend/internal/helper"
	"github.com/ricky7171/test_wa_backend/internal/usecase/chat"
	"github.com/ricky7171/test_wa_backend/internal/websocket"
)

type Chat struct {
	chatService chat.Service
}

func NewChatController(chatService chat.Service) *Chat {
	return &Chat{
		chatService: chatService,
	}
}

//get all chat with specific contact
func (ch *Chat) GetChat() gin.HandlerFunc {
	return func(c *gin.Context) {

		//1. get param contactId & lastId
		contactId := c.Param("contactId")
		lastId := c.Param("lastId")

		//2. validate input
		if contactId == "" || lastId == "" {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", "contact id or last id cannot be empty"))
			c.Abort()
			return
		}

		//3. get chat
		chats, err := ch.chatService.GetChatByContact(contactId, lastId)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", chats))

	}
}

//send new message to other user
func (ch *Chat) NewMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. read current user id, current user name, phone receiver, message, and contact id
		currentUserId := c.GetString("userId")
		currentUserName := c.GetString("name")

		var input struct {
			Phone   string `json:"phone"`
			Message string `json:"message"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//2. validate input
		if currentUserId == "" || currentUserName == "" || input.Phone == "" || input.Message == "" {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", "Phone or password cannot be empty"))
			c.Abort()
			return
		}

		//3. make prepared message
		message, err := ch.chatService.MakePreparedMessage(currentUserId, currentUserName, input.Phone, input.Message)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. send message object to broadcast channel
		websocket.MainHub.Broadcast <- *message

		//5. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", message))

	}
}
