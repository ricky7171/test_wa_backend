package helper

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SignedTokenDetails is representation of JWT Token payload
type SignedTokenDetails struct {
	Name  string
	Phone string
	ID    string
	jwt.StandardClaims
}

// SignedRefreshTokenDetails is representation of JWT Refresh Token payload
type SignedRefreshTokenDetails struct {
	ID string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

// GenerateAllTokens generates both the detailed token and refresh token
func GenerateAllTokens(name string, phone string, userId primitive.ObjectID) (signedToken string, signedRefreshToken string, err error) {
	userIdString := userId.Hex()

	//1. generate claims for token payload
	//token will expired 24 hours
	claims := &SignedTokenDetails{
		Name:  name,
		Phone: phone,
		ID:    userIdString,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	//2. generate refresh claims for refresh token payload
	//refresh token will expired 168 hours (1 week)
	refreshClaims := &SignedRefreshTokenDetails{
		ID: userIdString,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	//3. generate token and refresh token according to claims & refreshClaims
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		token = ""
		refreshToken = ""
		fmt.Println("error : ", err)
	}

	return token, refreshToken, err
}

//ValidateToken validates the jwt token
//convert token jadi data user
func ValidateToken(signedToken string) (claims *SignedTokenDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedTokenDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedTokenDetails)
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

//convert refresh_token to SignedRefreshTokenDetails that contain user id
func ValidateRefreshToken(signedToken string) (claims *SignedRefreshTokenDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedRefreshTokenDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedRefreshTokenDetails)
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
