package guard

import (
	"YenExpress/config"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type Identifier struct {
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	Audience  string `json:"audience"`
	UserId    uint   `json:"userId"`
	UserName  string `json:"userName,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Sex       string `json:"sex"`
}

func (data *Identifier) Valid() error {
	if data.Issuer != config.ServerDomain || data.Audience != config.WebClientDomain {
		return ErrInvalidToken
	}
	return nil
}

type Payload struct {
	UserId     uint      `json:"userId"`
	SessionID  string    `json:"sessionID"`
	Expiration time.Time `json:"expiration"`
	Issuer     string    `json:"issuer"`
	Class      string    `json:"class"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.Expiration) {
		return ErrExpiredToken
	} else if payload.Issuer != config.ServerDomain {
		return ErrInvalidToken
	}
	return nil
}

type JWTMaker struct {
	SecretKey string
}

var (
	ErrInvalidToken = errors.New("Invalid or unrecognized Token")
	ErrExpiredToken = errors.New("Token Expired")
)

func (maker *JWTMaker) CreateToken(payload Payload) (token string) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	signed_token, _ := jwtToken.SignedString([]byte(maker.SecretKey))
	return signed_token
}

func (maker *JWTMaker) CreateIdentifier(payload Identifier) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	signed_token, _ := jwtToken.SignedString([]byte(maker.SecretKey))
	return signed_token
}

func (maker *JWTMaker) VerifyToken(token, variety string) (*Patient, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.SecretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	var user *Patient
	err = config.DB.Where("ID = ?", payload.UserId).First(&user).Error
	if err != nil || payload.Class != variety {
		return nil, ErrInvalidToken
	}
	return user, err

}

func (maker *JWTMaker) GetPayloadFromToken(token, variety string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.SecretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

var Bearer = &JWTMaker{config.JwtSecret}
