package ratelimiter

import (
	"YenExpress/config"
)

var (
	PatientAPIKeyLimiter = FailRateLimiter{
		maxWrongAttemptsByIPperDay: 15, maxConsecutiveFails: 3,
		client: config.PatientRedisClient, TotalFailsKeyPrefix: "patient_apiKey_fail_per_day",
		ConsecutiveFailsKeyPrefix: "patient_apiKey_fail_consecutive",
	}

	PatientLoginLimiter = FailRateLimiter{
		maxWrongAttemptsByIPperDay: 100, maxConsecutiveFails: 10,
		client: config.PatientRedisClient, TotalFailsKeyPrefix: "patient_login_fail_ip_per_day",
		ConsecutiveFailsKeyPrefix: "patient_login_fail_consecutive_email_and_ip",
	}

	PatientCreateEmailOTPLimiter = FailRateLimiter{
		maxWrongAttemptsByIPperDay: 30, maxConsecutiveFails: 5,
		client: config.PatientRedisClient, TotalFailsKeyPrefix: "patient_create_email_otp_fail_per_day",
		ConsecutiveFailsKeyPrefix: "patient_create_email_otp_fail_consecutive",
	}

	PatientValidateEmailOTPLimiter = FailRateLimiter{
		maxWrongAttemptsByIPperDay: 20, maxConsecutiveFails: 5,
		client: config.PatientRedisClient, TotalFailsKeyPrefix: "patient_validate_email_otp_fail_per_day",
		ConsecutiveFailsKeyPrefix: "patient_validate_email_otp_fail_consecutive",
	}
)
