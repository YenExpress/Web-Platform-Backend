package config

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetIPAddress(c *gin.Context) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := c.Request.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := c.Request.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}

func GetBearerToken(c *gin.Context) (string, error) {
	credentials := c.Request.Header.Get("Authorization")
	credList := strings.Split(credentials, ",")
	if len(credList) > 2 {
		return "", errors.New("Authorization Header with incorrect format")
	}
	token := ""
	for _, val := range credList {
		if strings.HasPrefix(val, "Bearer ") {
			token = strings.TrimPrefix(val, "Bearer ")
		}
	}
	if token == "" {
		return "", errors.New("Bearer Unknown")
	}
	return token, nil
}

func GetAPIKey(c *gin.Context) (string, error) {
	credentials := c.Request.Header.Get("Authorization")
	credList := strings.Split(credentials, ",")
	if len(credList) > 2 {
		return "", errors.New("Authorization Header with incorrect format")
	}
	apiKey := ""
	for _, val := range credList {
		if strings.HasPrefix(val, "APIKey ") {
			apiKey = strings.TrimPrefix(val, "APIKey ")
		}
	}
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
