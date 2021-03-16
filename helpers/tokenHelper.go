package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"wa/database"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SignedDetails : ini adalah representasi dari token JWT
type SignedDetails struct {
	Name  string
	Phone string
	Uid   string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// GenerateAllTokens generates both the detailed token and refresh token
func GenerateAllTokens(name string, phone string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Name:  name,
		Phone: phone,
		Uid:   uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

//ValidateToken validates the jwt token
//convert token jadi data user
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}

//UpdateAllTokens renews the user tokens when they login
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	//1. buat context dengan timeout 100 detik
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//2. buat object BSON
	var updateObj primitive.D

	//3. isi object BSON :
	//{"token" : signedToken (DIAMBIL DARI PARAMETER), "refresh_token" : signedRefreshToken (DIAMBIL DARI PARAMETER), "updated_at" : (TIMENOW)}
	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	//4. bkin variabel upsert dimana ini digunakan sebagai penanda, bahwa
	//kalau datanya tidak ada maka insert data tersebut (sama persis kek upsertnya di laravel)
	upsert := true

	//5. buat object bson bernama filter lalu isi :
	//{"user_id" : userId (DIAMBIL DARI PARAMETER)}
	filter := bson.M{"user_id": userId}

	//6. buat object opt tipenya updateOptions dimana ini merepresentasikan option pada saat ngeupdate
	//apa optionnya ? ya upsert = &upsert yaitu upsert = true,
	//jadi tar kalau ngeupdate, kalau misal datanya gak ada, ya insert aja
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	//7. update user dengan id userId (diambil dari parameter) dengan data updateObj
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

	return
}