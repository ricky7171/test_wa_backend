package helper

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ricky7171/test_wa_backend/internal/failure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//interface Token
type Token interface {
	GenerateAllTokens(name string, phone string, userId primitive.ObjectID) (signedToken string, signedRefreshToken string, err error)
	ValidateToken(signedToken string) (claims *SignedTokenDetails, err error)
	ValidateRefreshToken(signedToken string) (claims *SignedRefreshTokenDetails, err error)
}

//implementation of Token interface (act as "class")
type TokenJWT struct {
}

// SignedTokenDetails is representation of JWT Token payload (act as data type)
type SignedTokenDetails struct {
	Name  string
	Phone string
	ID    string
	jwt.StandardClaims
}

// SignedRefreshTokenDetails is representation of JWT Refresh Token payload (act as data type)
type SignedRefreshTokenDetails struct {
	ID string
	jwt.StandardClaims
}

func NewTokenJwt() *TokenJWT {
	return &TokenJWT{}
}

// GenerateAllTokens function is used for generates both the token and refresh token
func (t *TokenJWT) GenerateAllTokens(name string, phone string, userId primitive.ObjectID) (signedToken string, signedRefreshToken string, err error) {
	secret := os.Getenv("SECRET_KEY")
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
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", "", failure.ErrGenerateToken()
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))
	if err != nil {
		return "", "", failure.ErrGenerateToken()
	}

	return token, refreshToken, err
}

//ValidateToken function used to validates the jwt token
func (t *TokenJWT) ValidateToken(signedToken string) (claims *SignedTokenDetails, err error) {
	secret := os.Getenv("SECRET_KEY")
	//1. convert claims to token
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedTokenDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, failure.ErrFailedParseClaim()
	}

	//2. check wether token is valid or invalid
	claims, ok := token.Claims.(*SignedTokenDetails)
	if !ok {
		return nil, failure.ErrTokenInvalid()
	}

	//3. check wether token was expired or not
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, failure.ErrTokenExpired()
	}

	return claims, nil
}

//validateRefreshToken function used to validates the jwt refresh token
func (t *TokenJWT) ValidateRefreshToken(signedToken string) (claims *SignedRefreshTokenDetails, err error) {
	secret := os.Getenv("SECRET_KEY")
	//1. convert claims to token
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedRefreshTokenDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, failure.ErrFailedParseClaim()
	}

	//2. check wether token is valid or invalid
	claims, ok := token.Claims.(*SignedRefreshTokenDetails)
	if !ok {
		return nil, failure.ErrTokenInvalid()
	}

	//3. check wether token was expired or not
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, failure.ErrTokenExpired()
	}

	return claims, nil
}
