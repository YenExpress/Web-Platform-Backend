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

func (maker *JWTMaker) CreatePatientToken(duration time.Time, user_id uint, variety string) (token string) {
	payload := &Payload{
		UserId: user_id, Class: variety,
		Expiration: duration,
		Issuer:     config.ServerDomain}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signed_token, _ := jwtToken.SignedString([]byte(maker.SecretKey))
	return signed_token
}

func (maker *JWTMaker) CreatePatientIdToken(userId uint, lastName, firstName, email, sex string) string {
	payload := &Identifier{UserId: userId, FirstName: firstName,
		LastName: lastName, Email: email, Sex: sex,
		Issuer: config.ServerDomain, Audience: config.WebClientDomain,
		Subject: "Token Bearing User Identity"}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signed_token, _ := jwtToken.SignedString([]byte(maker.SecretKey))
	return signed_token
}

func (maker *JWTMaker) VerifyPatientToken(token, variety string) (*Patient, error) {
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

var PatientJWTMaker = &JWTMaker{config.JwtSecret}

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
