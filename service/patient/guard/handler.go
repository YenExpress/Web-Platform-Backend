package guard

import (
	"YenExpress/config"
	"time"

	limiter "github.com/codeNino/ratelimiter"
)

func ValidateToken(token string, variety string) error {
	return Bearer.VerifyToken(token, variety)
}

func GetTokenPayload(token, variety string) (*Payload, error) {
	load, err := Bearer.GetPayloadFromToken(token, variety)
	if err == nil {

		return load, nil
	} else {
		return &Payload{}, err
	}

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
