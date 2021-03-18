package controllers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
	db "wa/database"
	helper "wa/helpers"
	"wa/hub"
	"wa/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	data := []byte(password)
	bytes := md5.Sum(data)
	if len(bytes) == 0 {
		log.Panic("cannot hash password")
	}

	return string(hex.EncodeToString(bytes[:]))
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
//userPassword is plain password
//providedPassword is hashed password
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	//compareHashAndPassword(hashed password, plain password)
	var check bool
	var msg string
	encryptedUserPassword := HashPassword(userPassword)
	if encryptedUserPassword == providedPassword {
		check = true
		msg = ""
	} else {
		check = false
		msg = fmt.Sprintf("login or passowrd is incorrect")
	}

	return check, msg
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. make ctx with timeout 100 second
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//2. make user model
		var user models.User

		//3. read JSON request then store to "user" variable (model)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//4. validate user
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//5. check wether phone is already in the database or not
		count, err := db.UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel() //defer cancel() used to clean go routine memory after this function is done
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this phone already exists"})
			return
		}

		//6. hash user password
		password := HashPassword(*user.Password)
		user.Password = &password

		//7. fill attribute : created_at, updated_at, id, token, and refresh_token
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		//8. generate token.
		//generate JWT token from name, phone, and user id
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Name, *user.Phone, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		//9. insert user to database
		resultInsertionNumber, insertErr := db.UserCollection.InsertOne(ctx, user)
		defer cancel()
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//10. send response to client
		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//4. search user that have client request phone number, then save that user to "founduser" variable
		err := db.UserCollection.FindOne(ctx, bson.M{"phone": user.Phone}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
			return
		}

		//5. check userfound password with client request password
		//*user.Password is password plain
		//*foundUser.Password is hashed password
		//if wrong, return error : login or passowrd is incorrect
		//this process will take long time (about 1 second), because bcrypt is complex
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//6. generate token.
		//generate token from name, phone, and user id
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Name, *foundUser.Phone, *&foundUser.User_id)

		//7. update user with new token and new refresToken
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		//8. remove password attribute, because it will send to client
		foundUser.Password = nil

		//9. send response to client
		c.JSON(http.StatusOK, foundUser)

	}
}

func ConnectWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		hub.ServeWs(c.Writer, c.Request, userID)
	}
}

//get all chat with specific room
func GetChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. get param room_id & last_id
		roomID := c.Param("room_id")
		roomObjectID, err := primitive.ObjectIDFromHex(roomID)
		lastID := c.Param("last_id")
		var lastObjectID primitive.ObjectID
		if lastID != "" {
			lastObjectID, err = primitive.ObjectIDFromHex(lastID)
		}

		//2. make context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//3. get data room with all chat in that room
		matchRoomIDPipeline := bson.D{{"$match", bson.D{{"_id", roomObjectID}}}}
		unwindPipeline := bson.D{{"$unwind", "$chat_history"}}
		sortPipeline := bson.D{{"$sort", bson.D{{"chat_history._id", -1}}}}
		paginatePipeline := bson.D{{"$match", bson.D{{"chat_history._id", bson.D{{"$lt", lastObjectID}}}}}}
		limitPipeline := bson.D{{"$limit", 20}}
		groupPipeline := bson.D{{"$group", bson.D{{"_id", "$_id"}, {"chat_history", bson.M{"$push": "$chat_history"}}}}}

		var cursor *mongo.Cursor

		if lastID == "nil" { //means this client request first page

			cursor, err = db.RoomCollection.Aggregate(
				ctx,
				mongo.Pipeline{
					matchRoomIDPipeline,
					unwindPipeline,
					sortPipeline,
					limitPipeline,
					groupPipeline,
				},
			)
		} else { //means this client request NOT the first page

			cursor, err = db.RoomCollection.Aggregate(
				ctx,
				mongo.Pipeline{
					matchRoomIDPipeline,
					unwindPipeline,
					sortPipeline,
					paginatePipeline,
					limitPipeline,
					groupPipeline,
				},
			)
		}
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var result []bson.M
		if err = cursor.All(ctx, &result); err != nil {
			log.Fatal(err)
		}

		//5. send response to client
		c.JSON(http.StatusOK, result)

	}
}

//get contact
//contact is peoples who have interacted before with specific user id
func GetContact() gin.HandlerFunc {
	return func(c *gin.Context) {

		//1. get param user_id and make objectID from that
		userID := c.GetString("user_id")
		userObjetID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error convert from string to objectID"})
		}

		//2. make context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//3. get all data chat WHERE it contains userId
		cursor, err := db.RoomCollection.Aggregate(
			ctx,
			mongo.Pipeline{
				bson.D{
					{
						"$match", bson.M{
							"users": bson.M{
								"$all": []interface{}{userObjetID},
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

		//4. store all result in allContacts
		var allContacts []bson.M
		if err = cursor.All(ctx, &allContacts); err != nil {
			log.Fatal(err)
		}

		//5. send response to client
		c.JSON(http.StatusOK, allContacts)

	}
}

//send new message to other user
func NewMessage() gin.HandlerFunc {
	return func(c *gin.Context) {

		//1. make ctx with timeout 100 second
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//2. make chat model that store client request
		var newChat models.NewChat

		//3. read request from client and store to "newChat" variable
		if err := c.BindJSON(&newChat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//4. make object model message
		m := models.Message{Data: newChat.Message, FromUserId: c.GetString("user_id")}

		//5. if client doesn't send room_id and send phone number, then it need to search wether client have interacted with this phone number before or not
		if newChat.Room_id == "" && newChat.Phone != "" {
			//5.1. make user model
			var user models.User

			//5.b. get user data according that phone number
			err := db.UserCollection.FindOne(ctx, bson.M{"phone": newChat.Phone}).Decode(&user)
			defer cancel()

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			//5.c. add user_id to message object
			m.ToUserId = user.User_id
		} else if newChat.Room_id != "" { //if client send room_id
			//5.a. add rooom_id to message object
			m.Room_id = newChat.Room_id
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "room id and phone not found"})
			return
		}

		//6. send message object to broadcast channel
		hub.MainHub.Broadcast <- m

		//7. send response to client
		c.JSON(http.StatusOK, map[string]bool{
			"success": true,
		})

	}
}
