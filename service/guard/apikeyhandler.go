package guard

import (
	"YenExpress/config"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	limiter "github.com/codeNino/ratelimiter"

	"github.com/gin-gonic/gin"
)

// func GetAPIKey(c *gin.Context) (string, error) {
// 	credentials := c.Request.Header.Get("Authorization")
// 	credList := strings.Split(credentials, ",")
// 	if len(credList) > 2 {
// 		return "", errors.New("Authorization Header with incorrect format")
// 	}
// 	apiKey := ""
// 	for _, val := range credList {
// 		if strings.HasPrefix(val, "APIKey ") {
// 			apiKey = strings.TrimPrefix(val, "APIKey ")
// 		}
// 	}
// 	if apiKey == "" {
// 		return "", errors.New("APIKey Unknown")
// 	}
// 	return apiKey, nil
// }

func GetAPIKey(c *gin.Context) (string, error) {
	apiKey := c.Request.Header.Get("x-api-key")
	if apiKey == "" {
		return "", errors.New("APIKey Unknown")
	}
	return apiKey, nil
}

func ValidateAPIKey(providedKey, validKey string) (bool, error) {

	decodedProvided, err := hex.DecodeString(providedKey)
	if err != nil {
		return false, errors.New("Invalid API key")
	}
	decodedValid, _ := hex.DecodeString(validKey)
	if ok := apiKeyIsValid(decodedProvided, decodedValid); !ok {
		log.Println("No Matching API key found")
		return false, errors.New("Invalid API key")
	}
	return true, nil

}

func apiKeyIsValid(provided, valid []byte) bool {
	contentEqual := subtle.ConstantTimeCompare(valid, provided) == 1

	if contentEqual {
		return true
	}
	return false
}

func HashAPIKey(rawKey string) string {
	hash := sha256.Sum256([]byte(rawKey))
	return fmt.Sprintf("%x", hash)
}

var APIKeyLimiter = limiter.RateLimiter{
	TotalLimit: 15, BurstLimit: 3, MaxTime: time.Hour * 24, BurstPeriod: time.Minute * 30,
	Client: config.PatientRedisClient, TotalLimitPrefix: "apiKey_fail_per_day",
	BurstLimitPrefix: "apiKey_fail_consecutive",
}
