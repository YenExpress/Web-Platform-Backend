package guard

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken = errors.New("Invalid or unrecognized Token")
	ErrExpiredToken = errors.New("Token Expired")
)

type JWTStrategy struct {
	SecretKey string
	UserValid func(id interface{}) bool
}

func (maker *JWTStrategy) CreateToken(payload Bearer) (token string) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	signed_token, _ := jwtToken.SignedString([]byte(maker.SecretKey))
	return signed_token
}

func (maker *JWTStrategy) CreateIdentifier(payload Identifier) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	signed_token, _ := jwtToken.SignedString([]byte(maker.SecretKey))
	return signed_token
}

func (maker *JWTStrategy) VerifyToken(token, variety string) error {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.SecretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Bearer{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return ErrExpiredToken
		}
		return ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Bearer)
	if !ok {
		return ErrInvalidToken
	}
	userExists := maker.UserValid(payload.UserId)
	// var user *Patient
	// err = config.DB.Where("ID = ?", payload.UserId).First(&user).Error
	if !userExists || payload.Class != variety {
		return ErrInvalidToken
	}
	return err

}

func (maker *JWTStrategy) GetPayloadFromToken(token, variety string) (*Bearer, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.SecretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Bearer{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Bearer)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

// func GetBearerToken(c *gin.Context) (string, error) {
// 	credentials := c.Request.Header.Get("Authorization")
// 	credList := strings.Split(credentials, ",")
// 	if len(credList) > 2 {
// 		return "", errors.New("Authorization Header with incorrect format")
// 	}
// 	token := ""
// 	for _, val := range credList {
// 		if strings.HasPrefix(val, "Bearer ") {
// 			token = strings.TrimPrefix(val, "Bearer ")
// 		}
// 	}
// 	if token == "" {
// 		return "", errors.New("Bearer Unknown")
// 	}
// 	return token, nil
// }

func GetBearerToken(c *gin.Context) (string, error) {
	authString := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(authString, "Bearer ") {
		token := strings.TrimPrefix(authString, "Bearer ")
		if token == "" {
			return "", errors.New("Bearer Unknown")
		}
		return token, nil

	}
	return "", errors.New("Bearer Unknown")
}
