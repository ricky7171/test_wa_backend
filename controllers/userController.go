package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"wa/database"
	helper "wa/helpers"
	"wa/hub"
	"wa/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var roomCollection *mongo.Collection = database.OpenCollection(database.Client, "rooms")
var validate = validator.New()

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
//userPassword adalah plain password
//providedPassword adalah hashed password
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	//compareHashAndPassword(hashed password, plain password)
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

//CreateUser is the api used to tget a single user
func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. buat ctx dengan timeout 100 detik
		//dimana ctx ini bakal digunakan saat nge query
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//2. siapkan user model
		var user models.User

		//3. baca JSON request dan masukan ke variabel user (model)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//4. validate user tsb (kek di laravel)
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//5. ngecek apakah phonenya sudah ada atau belum
		count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel() //defer cancel() gunanya untuk membersihkan memory yang ada sangkut pautnya dengan ctx
		//nah kenapa setiap nge query harus di defer cancel(), karena saat ngequery kita mengirimkan ctx
		//jadi kalau 1x pake ctx, ya di defer cancel() nya sekali, kalau 2x pake ctx, ya berarti 2x harus di defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this phone already exists"})
			return
		}

		//6. hash passwordnya
		password := HashPassword(*user.Password)
		user.Password = &password

		//7. mengisi attribute : created_at, updated_at, id, token, dan refresh_token
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		//8. generate token.
		//generate token dari name, phone, dan user id
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Name, *user.Phone, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		//9. insert ke database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		defer cancel()
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//10. kirim response
		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

//Login is the api used to tget a single user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. buat ctx dengan timeout 100 detik
		//dimana ctx ini bakal digunakan saat nge query
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//2. siapkan user model untuk variabel user dan foundUser
		//dimana user untuk request dan foundUser untuk hasil query
		var user models.User
		var foundUser models.User

		//3. baca JSON request dan masukan ke variabel user (model)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//4. cari data user dengan phone sesuai dengan phone request
		//lalu setelah ketemu kan bentuknya masih JSON. trus di decode ke variabel foundUser
		err := userCollection.FindOne(ctx, bson.M{"phone": user.Phone}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
			return
		}

		//5. cek password dari phone itu sudah sesuai dengan password yang dikirim oleh user atau belum
		//*user.Passowrd itu adalah password plain
		//*foundUser.Password itu adalah hashed password
		//kalau salah, langsung balikin error : login or passowrd is incorrect
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		//6. generate token.
		//generate token dari name, phone, dan user id
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Name, *foundUser.Phone, *&foundUser.User_id)

		//7. update user tersebut dengan token, refresToken yang baru
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		//8. hapus attribute password
		foundUser.Password = nil

		//9. kirim response
		c.JSON(http.StatusOK, foundUser)

	}
}

//get all chat with specific room
func GetChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. get param room_id
		roomID := c.Param("room_id")

		//2. buat context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//3. siapkan model room
		var room models.Room

		//4. nge get data chat dari room tersebut
		err := roomCollection.FindOne(ctx, bson.M{"room_id": roomID}).Decode(&room)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, room)

	}
}

//get contact
func GetContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. get param user_id
		userID := c.GetString("user_id")
		userObjetID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error convert from string to objectID"})
		}

		//2. buat context
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//3. nge get data chat dari room tersebut
		cursor, err := roomCollection.Aggregate(
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

		fmt.Println(cursor)

		var allContacts []bson.M
		if err = cursor.All(ctx, &allContacts); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allContacts)

	}
}

//new message
func NewMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1. buat ctx dengan timeout 100 detik
		//dimana ctx ini bakal digunakan saat nge query
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//2. siapkan new chat model
		//dimana new chat ini akan menampung request dari user
		var newChat models.NewChat

		//3. baca JSON request dan masukan ke variabel newChat (model)
		if err := c.BindJSON(&newChat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//4. membuat object model message
		m := models.Message{Data: newChat.Message, FromUserId: c.GetString("user_id")}

		//saat user mengirimkan room_id ke API ini, maka ya tidak usah
		//melakukan query untuk mendapatkan informasi penerima pesan ini

		//5. tapi kalau ternyata user tidak mengirimmkan room_id tapi ngirim phone, berarti kita harus cari informasi penerimanya lewat phone. Informasi tsb digunakan untuk nge get room idnya.
		if newChat.Room_id == "" && newChat.Phone != "" {
			//5.a. siapkan model user
			var user models.User

			//5.b. nge get data user sesuai phone penerima yang dikirim
			err := userCollection.FindOne(ctx, bson.M{"phone": newChat.Phone}).Decode(&user)
			defer cancel()

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			//5.c. tambahkan id penerima ke object message
			m.ToUserId = user.User_id
		} else if newChat.Room_id != "" { //jika user mengirimkan room_id
			//5.a. tambahkan id room ke object message
			m.Room_id = newChat.Room_id
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "room id and phone not found"})
			return
		}

		//6. kirimkan object message ke channel boradcast
		hub.MainHub.Broadcast <- m

		//7. kirim response
		c.JSON(http.StatusOK, map[string]bool{
			"success": true,
		})

	}
}
