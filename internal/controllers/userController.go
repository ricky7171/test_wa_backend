package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	helper "github.com/ricky7171/test_wa_backend/internal/helpers"
	"github.com/ricky7171/test_wa_backend/internal/hub"
	"github.com/ricky7171/test_wa_backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	return string(bytes), err
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
//userPassword is plain password
//providedPassword is hashed password
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

func Register(dbInstance *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. make ctx with timeout 100 second
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//2. make user model
		var user models.User

		//3. read JSON request then store to "user" variable (model)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. validate user
		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", validationErr.Error()))
			c.Abort()
			return
		}

		//5. check wether phone is already in the database or not
		count, err := dbInstance.Collection("users").CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel() //defer cancel() used to clean go routine memory after this function is done
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "error occured while checking for the phone"))
			c.Abort()
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "this phone already exists"))
			c.Abort()
			return
		}

		//6. hash user password
		password, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "failed to hash password"))
			c.Abort()
			return
		}
		user.Password = password

		//7. fill attribute : createdAt, updatedAt, id, token, and refreshToken
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.ID = primitive.NewObjectID()

		//8. insert user to database
		resultInsert, insertErr := dbInstance.Collection("users").InsertOne(ctx, user)
		defer cancel()
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "User item was not created"))
			c.Abort()
			return
		}

		//9. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", resultInsert))

	}
}

func Login(dbInstance *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. make ctx with timeout 100 second
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//2. make user model for "user" variable and "foundUser" variable
		//"user" variable used to store request
		//"founduser" variable used to store query result
		var user models.User
		var foundUser models.User

		//3. read request and store to "user" variable
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. validate request
		validationErr := validate.StructPartial(user, "Phone", "Password")

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", validationErr.Error()))
			c.Abort()
			return
		}

		//5. search user that have client request phone number, then save that user to "founduser" variable
		err := dbInstance.Collection("users").FindOne(ctx, bson.M{"phone": user.Phone}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "login or password is incorrect"))
			c.Abort()
			return
		}

		//6. check userfound password with client request password
		//user.Password is password plain
		//foundUser.Password is hashed password
		//if wrong, return error : login or passowrd is incorrect
		//this process will take long time (about 1 second), because bcrypt is complex
		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", msg))
			c.Abort()
			return
		}

		//7. generate token.
		//generate token from name, phone, and user id
		token, refreshToken, _ := helper.GenerateAllTokens(foundUser.Name, foundUser.Phone, foundUser.ID)

		//8. update user with new token and new refresToken
		foundUser.Token = token
		foundUser.RefreshToken = refreshToken

		//9. remove password attribute, because it will send to client
		foundUser.Password = ""

		//10. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", foundUser))

	}
}

func ConnectWs() gin.HandlerFunc {

	return func(c *gin.Context) {
		userID := c.Param("userId")
		hub.ServeWs(c.Writer, c.Request, userID)
	}
}

//get all chat with specific contact
func GetChat(dbInstance *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. get param contactId & lastId
		contactID := c.Param("contactId")
		contactObjectID, err := primitive.ObjectIDFromHex(contactID)
		lastID := c.Param("lastId")
		var lastObjectID primitive.ObjectID
		if lastID != "" {
			lastObjectID, _ = primitive.ObjectIDFromHex(lastID)
		}
		//2. make context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//3. get data chat according to that contact id
		matchContactIDPipeline := bson.D{{"$match", bson.D{{"contactId", contactObjectID}}}}
		sortPipeline := bson.D{{"$sort", bson.D{{"_id", -1}}}}
		paginatePipeline := bson.D{{"$match", bson.D{{"_id", bson.D{{"$lt", lastObjectID}}}}}}
		limitPipeline := bson.D{{"$limit", 20}}
		groupPipeline := bson.D{{"$project", bson.D{{"_id", 1}, {"contactId", 1}, {"senderId", 1}, {"message", 1}, {"createdAt", 1}}}}

		var cursor *mongo.Cursor

		if lastID == "nil" { //means this client request first page

			cursor, err = dbInstance.Collection("chats").Aggregate(
				ctx,
				mongo.Pipeline{
					matchContactIDPipeline,
					sortPipeline,
					limitPipeline,
					groupPipeline,
				},
			)
		} else { //means this client request NOT the first page

			cursor, err = dbInstance.Collection("chats").Aggregate(
				ctx,
				mongo.Pipeline{
					matchContactIDPipeline,
					sortPipeline,
					paginatePipeline,
					limitPipeline,
					groupPipeline,
				},
			)
		}

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", err))
			c.Abort()
			return
		}

		var chats []models.Chat
		if err = cursor.All(ctx, &chats); err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "Invalid data format from DB"))
			c.Abort()
			return
		}

		//5. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", chats))

	}
}

//get contact
//contact is peoples who have interacted before with specific user id
func GetContact(dbInstance *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		//1. get param userId and make objectID from that
		userID := c.GetString("userId")
		userObjetID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "Error convert from string to objectID"))
			c.Abort()
			return
		}

		//2. make context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//3. get all data chat WHERE it contains userId
		cursor, err := dbInstance.Collection("contacts").Aggregate(
			ctx,
			mongo.Pipeline{
				bson.D{
					{
						"$match", bson.M{
							"users": bson.M{
								"$in": []interface{}{userObjetID},
							},
						},
					},
				},
				bson.D{
					{
						"$lookup", bson.M{
							"from":         "users",
							"localField":   "users",
							"foreignField": "_id",
							"as":           "users_info",
						},
					},
				},
				bson.D{
					{
						"$project", bson.M{
							"users_info._id":  1,
							"users_info.name": 1,
						},
					},
				},
			},
		)

		defer cancel()

		//4. convert cursor to bson.M
		var allContacts []bson.M
		if err = cursor.All(ctx, &allContacts); err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "Data is invalid"))
			c.Abort()
			return
		}

		//5. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", allContacts))

	}
}

//send new message to other user
func NewMessage(dbInstance *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		//1. make ctx with timeout 100 second
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//2. make chat model that store client request
		var newChat models.NewChat

		//3. read request from client and store to "newChat" variable
		if err := c.BindJSON(&newChat); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. make object model message
		m := models.Message{Data: newChat.Message, FromUserId: c.GetString("userId")}

		//5. if client doesn't send contactId but send phone number, then it need to search wether client have contact with this phone number or not
		if newChat.ContactId == "" && newChat.Phone != "" {
			//5.1. make user model
			var user models.User

			//5.b. get user data according that phone number
			err := dbInstance.Collection("users").FindOne(ctx, bson.M{"phone": newChat.Phone}).Decode(&user)
			defer cancel()

			if err != nil {
				c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", err.Error()))
				c.Abort()
				return
			}
			//5.c. add userId to message object
			m.ToUserId = user.ID.Hex()
		} else if newChat.ContactId != "" { //if client send contatId
			//5.a. add rooomId to message object
			m.ContactId = newChat.ContactId
		} else {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "Contact id and phone not found"))
			c.Abort()
			return
		}

		//6. send message object to broadcast channel
		hub.MainHub.Broadcast <- m

		//7. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", m))

	}
}

//token refresh
//used to refresh access token that has been expired
func RefreshToken(dbInstance *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. make ctx with timeout 100 second
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//2. get token from body
		var request map[string]interface{}

		//3. read request from client and store to "request" variable
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", err.Error()))
			c.Abort()
			return
		}

		//4. check if token is present
		plainToken, ok := request["refresh_token"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, helper.FormatResponse("error", "refresh_token is not present"))
			c.Abort()
			return
		}

		//5. change plain token to be signedDetails that contains user id
		claims, errMessage := helper.ValidateRefreshToken(plainToken)
		if errMessage != "" {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", errMessage))
			c.Abort()
			return
		}

		//6. get user with ID that get from claims
		var user models.User
		userID, _ := primitive.ObjectIDFromHex(claims.ID)
		err := dbInstance.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.FormatResponse("error", "user not found"))
			c.Abort()
			return
		}

		//7. generate new token and refresh token
		token, refreshToken, _ := helper.GenerateAllTokens(user.Name, user.Phone, user.ID)

		//8. update user with new token and new refresToken
		user.Token = token
		user.RefreshToken = refreshToken

		//9. remove password attribute, because it will send to client
		user.Password = ""

		//10. send response to client
		c.JSON(http.StatusOK, helper.FormatResponse("success", user))
	}
}
