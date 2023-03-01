package guard

import (
	"YenExpress/config"
	limiter "YenExpress/ratelimiter"
	"strings"

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
	APIKeyLimiter = limiter.FailRateLimiter{
		MaxWrongAttemptsByIPperDay: 15, MaxConsecutiveFails: 3,
		Client: config.PatientRedisClient, TotalFailsKeyPrefix: "patient_apiKey_fail_per_day",
		ConsecutiveFailsKeyPrefix: "patient_apiKey_fail_consecutive",
	}

	LoginLimiter = limiter.FailRateLimiter{
		MaxWrongAttemptsByIPperDay: 100, MaxConsecutiveFails: 10,
		Client: config.PatientRedisClient, TotalFailsKeyPrefix: "patient_login_fail_ip_per_day",
		ConsecutiveFailsKeyPrefix: "patient_login_fail_consecutive_email_and_ip",
	}

	CreateEmailOTPLimiter = limiter.FailRateLimiter{
		MaxWrongAttemptsByIPperDay: 30, MaxConsecutiveFails: 5,
		Client: config.PatientRedisClient, TotalFailsKeyPrefix: "patient_create_email_otp_fail_per_day",
		ConsecutiveFailsKeyPrefix: "patient_create_email_otp_fail_consecutive",
	}

	EmailValidationLimiter = limiter.FailRateLimiter{
		MaxWrongAttemptsByIPperDay: 20, MaxConsecutiveFails: 5,
		Client: config.PatientRedisClient, TotalFailsKeyPrefix: "patient_validate_email_otp_fail_per_day",
		ConsecutiveFailsKeyPrefix: "patient_validate_email_otp_fail_consecutive",
	}
)
