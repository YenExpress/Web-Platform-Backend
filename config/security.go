package config

// import (
// 	"YenExpress/responses"
// 	"crypto/sha256"
// 	"crypto/subtle"
// 	"encoding/hex"
// 	"errors"
// 	"net"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"

// 	"golang.org/x/crypto/bcrypt"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/google/uuid"
// )

// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

// func CheckPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	return err == nil
// }

// type Payload struct {
// 	TokenID   string    `json:"TokenID"`
// 	UserID    string    `json:"UserID"`
// 	ExpiredAt time.Time `json:"ExpiredAt"`
// }

// type Maker interface {
// 	CreateAccessToken(user_id string) (token string)
// 	VerifyAccessToken(token string) (*Payload, error)
// 	CreateRefreshToken(user_id string) (token string)
// 	VerifyRefreshToken(token string) (*Payload, error)
// }

// type JWTMaker struct {
// 	secretKey string
// }

// var (
// 	ErrInvalidToken = errors.New("token is invalid")
// 	ErrExpiredToken = errors.New("token has expired")
// )

// func (payload *Payload) Valid() error {
// 	if time.Now().After(payload.ExpiredAt) {
// 		return ErrExpiredToken
// 	}
// 	return nil
// }

// func (maker *JWTMaker) CreateAccessToken(user_id string) (token string) {
// 	payload, _ := NewAccessPayload(user_id)
// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
// 	signed_token, _ := jwtToken.SignedString([]byte(maker.secretKey))
// 	return signed_token
// }

// func (maker *JWTMaker) CreateRefreshToken(user_id string) (token string) {
// 	payload, _ := NewRefreshPayload(user_id)
// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
// 	signed_token, _ := jwtToken.SignedString([]byte(maker.secretKey))
// 	return signed_token
// }

// func (maker *JWTMaker) VerifyAccessToken(token string) (*Payload, error) {
// 	keyFunc := func(token *jwt.Token) (interface{}, error) {
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, ErrInvalidToken
// 		}
// 		return []byte(maker.secretKey), nil
// 	}
// 	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
// 	if err != nil {
// 		verr, ok := err.(*jwt.ValidationError)
// 		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
// 			return nil, ErrExpiredToken
// 		}
// 		return nil, ErrInvalidToken
// 	}
// 	payload, ok := jwtToken.Claims.(*Payload)
// 	if !ok {
// 		return nil, ErrInvalidToken
// 	}

// 	return payload, nil
// }

// func (maker *JWTMaker) VerifyRefreshToken(token string) (*Payload, error) {
// 	keyFunc := func(token *jwt.Token) (interface{}, error) {
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, ErrInvalidToken
// 		}
// 		return []byte(maker.secretKey), nil
// 	}
// 	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
// 	if err != nil {
// 		verr, ok := err.(*jwt.ValidationError)
// 		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
// 			return nil, ErrExpiredToken
// 		}
// 		return nil, ErrInvalidToken
// 	}
// 	payload, ok := jwtToken.Claims.(*Payload)
// 	if !ok {
// 		return nil, ErrInvalidToken
// 	}

// 	return payload, nil
// }

// func NewAccessPayload(user_id string) (*Payload, error) {
// 	tokenID := (uuid.New()).String()
// 	payload := &Payload{
// 		TokenID:   tokenID,
// 		UserID:    user_id,
// 		ExpiredAt: time.Now().Add(time.Hour * 24),
// 	}
// 	return payload, nil
// }

// func NewRefreshPayload(user_id string) (*Payload, error) {
// 	tokenID := (uuid.New()).String()
// 	payload := &Payload{
// 		TokenID:   tokenID,
// 		UserID:    user_id,
// 		ExpiredAt: time.Now().Add(time.Hour * 24 * 30),
// 	}
// 	return payload, nil
// }

// func NewJWTMaker() (Maker, error) {
// 	return &JWTMaker{SECRET_KEY}, nil
// }

// func GetAccessToken(user_id string) string {
// 	maker, _ := NewJWTMaker()
// 	return maker.CreateAccessToken(user_id)
// }

// func GetRefreshToken(user_id string) string {
// 	maker, _ := NewJWTMaker()
// 	return maker.CreateRefreshToken(user_id)
// }

// // decode payload and confirm token and  user validity
// func ValidateAccessToken(token string) (*Payload, error) {
// 	maker, _ := NewJWTMaker()
// 	payload, err := maker.VerifyAccessToken(token)
// 	if err != nil {
// 		return nil, err
// 	}
// 	activesessiondata, _ := UserSessionManager.Get_SessionData(payload.UserID)
// 	if activesessiondata.TokenID != payload.TokenID {
// 		return nil, ErrInvalidToken
// 	}
// 	return payload, err

// }

// func ValidateRefreshToken(token string) (*Payload, error) {
// 	maker, _ := NewJWTMaker()
// 	payload, err := maker.VerifyRefreshToken(token)
// 	if err != nil {
// 		return nil, err
// 	}
// 	activesessiondata, _ := UserSessionManager.Get_SessionData(payload.UserID)
// 	if activesessiondata.TokenID != payload.TokenID {
// 		return nil, ErrInvalidToken
// 	}
// 	return payload, err

// }

// // validator to check if old token was recently issued on server before granting access to refresh token API
// func AuthorizeRefreshToken(c *gin.Context) bool {
// 	status := false
// 	token := c.Request.Header.Get("Authorization")
// 	if strings.HasPrefix(token, "Bearer ") {
// 		token = strings.TrimPrefix(token, "Bearer ")
// 		valid_payload, err := ValidateAccessToken(token)
// 		if err != nil {
// 			c.JSON(http.StatusForbidden, responses.DefaultResponse{StatusCode: 401, Message: err.Error(), Status: false})
// 		} else if valid_payload != nil {
// 			status = false
// 		}
// 	} else {
// 		c.JSON(http.StatusUnauthorized, responses.DefaultResponse{StatusCode: 401, Message: "Unauthorized", Status: false})
// 	}
// 	return status
// }

// func GetUserFromToken(c *gin.Context) (id string) {
// 	token := c.Request.Header.Get("Authorization")
// 	if strings.HasPrefix(token, "Bearer ") {
// 		token = strings.TrimPrefix(token, "Bearer ")
// 		decoded_payload, err := ValidateAccessToken(token)
// 		if err == nil {
// 			return decoded_payload.UserID
// 		} else {
// 			return ""
// 		}
// 	}
// 	return ""
// }

// func AuthorizeEndpointAccess(c *gin.Context) bool {
// 	status := false
// 	token := c.Request.Header.Get("Authorization")
// 	if strings.HasPrefix(token, "Bearer ") {
// 		token = strings.TrimPrefix(token, "Bearer ")
// 		payload, err := ValidateAccessToken(token)
// 		if err != nil {
// 			c.JSON(http.StatusForbidden, responses.DefaultResponse{StatusCode: 401, Message: err.Error(), Status: false})
// 			status = false
// 		} else {
// 			if payload.ExpiredAt.Before(time.Now()) {
// 				c.JSON(http.StatusForbidden, responses.DefaultResponse{StatusCode: 403, Message: "Forbidden! Token Expired!!", Status: false})
// 			} else {
// 				status = true
// 			}
// 		}
// 	} else {
// 		c.JSON(http.StatusUnauthorized, responses.DefaultResponse{StatusCode: 401, Message: "Unauthorized!! Bearer Unknown", Status: false})
// 		status = false
// 	}
// 	return status
// }

// func ApiKeyMiddleware(cfg conf.Config, logger logging.Logger) (func(handler http.Handler) http.Handler, error) {
// 	apiKeyHeader := cfg.APIKeyHeader
// 	apiKeys := cfg.APIKeys
// 	apiKeyMaxLen := cfg.APIKeyMaxLen

// 	decodedAPIKeys := make(map[string][]byte)
// 	for name, value := range apiKeys {
// 		decodedKey, err := hex.DecodeString(value)
// 		if err != nil {
// 			return nil, err
// 		}

// 		decodedAPIKeys[name] = decodedKey
// 	}

// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			ctx := r.Context()

// 			apiKey, err := bearerToken(r, apiKeyHeader)
// 			if err != nil {
// 				logger.Errorw("request failed API key authentication", "error", err)
// 				RespondError(w, http.StatusUnauthorized, "invalid API key")
// 				return
// 			}

// 			if _, ok := apiKeyIsValid(apiKey, decodedAPIKeys); !ok {
// 				hostIP, _, err := net.SplitHostPort(r.RemoteAddr)
// 				if err != nil {
// 					logger.Errorw("failed to parse remote address", "error", err)
// 					hostIP = r.RemoteAddr
// 				}
// 				logger.Errorw("no matching API key found", "remoteIP", hostIP)

// 				RespondError(w, http.StatusUnauthorized, "invalid api key")
// 				return
// 			}

// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}, nil
// }

// // apiKeyIsValid checks if the given API key is valid and returns the principal if it is.
// func apiKeyIsValid(rawKey string, availableKeys map[string][]byte) (string, bool) {
// 	hash := sha256.Sum256([]byte(rawKey))
// 	key := hash[:]

// 	for name, value := range availableKeys {
// 		contentEqual := subtle.ConstantTimeCompare(value, key) == 1

// 		if contentEqual {
// 			return name, true
// 		}
// 	}

// 	return "", false
// }

// // bearerToken function omitted...
