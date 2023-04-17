package guard

import (
	"YenExpress/config"
	"strings"
	"time"

	limiter "github.com/codeNino/ratelimiter"

	"github.com/gin-gonic/gin"
)

func ValidateAccessToken(token string) (*Patient, error) {
	user, err := Bearer.VerifyToken(token, "access_token")
	if err != nil {
		return nil, err
	}
	return user, err
}

func ValidateRefreshToken(token string) error {
	_, err := Bearer.VerifyToken(token, "refresh_token")
	return err

}

func GetPatientIDFromToken(c *gin.Context) (*Patient, error) {
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		user, err := ValidateAccessToken(token)
		if err == nil {

			return user, nil
		} else {
			return &Patient{}, err
		}
	}
	return &Patient{}, ErrInvalidToken
}

var (
	APIKeyLimiter = limiter.RateLimiter{
		TotalLimit: 15, BurstLimit: 3, MaxTime: time.Hour * 24, BurstPeriod: time.Minute * 30,
		Client: config.PatientRedisClient, TotalLimitPrefix: "patient_apiKey_fail_per_day",
		BurstLimitPrefix: "patient_apiKey_fail_consecutive",
	}

	LoginLimiter = limiter.RateLimiter{
		TotalLimit: 100, BurstLimit: 10, MaxTime: time.Hour * 24, BurstPeriod: time.Hour * 1,
		Client: config.PatientRedisClient, TotalLimitPrefix: "patient_login_fail_ip_per_day",
		BurstLimitPrefix: "patient_login_fail_consecutive_email_and_ip",
	}

	CreateOTPLimiter = limiter.RateLimiter{
		TotalLimit: 30, BurstLimit: 5, MaxTime: time.Hour * 24, BurstPeriod: time.Minute * 30,
		Client: config.PatientRedisClient, TotalLimitPrefix: "patient_create_otp_fail_per_day",
		BurstLimitPrefix: "patient_create_otp_fail_consecutive",
	}

	EmailValidationLimiter = limiter.RateLimiter{
		TotalLimit: 30, BurstLimit: 5, MaxTime: time.Hour * 24, BurstPeriod: time.Minute * 30,
		Client: config.PatientRedisClient, TotalLimitPrefix: "patient_validate_email_otp_fail_per_day",
		BurstLimitPrefix: "patient_validate_email_otp_fail_consecutive",
	}
)
