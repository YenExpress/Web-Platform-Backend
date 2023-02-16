package ratelimiter

import (
	"YenExpress/config"
)

var PatientLoginLimiter = RateLimiter{
	maxWrongAttemptsByIPperDay: 100, maxConsecutiveFailsByEmailAndIP: 10,
	client: config.PatientRedisClient, SlowBruteKeyPrefix: "patient_login_fail_ip_per_day",
	ConsecutiveFailsKeyPrefix: "patient_login_fail_consecutive_email_and_ip",
}

var PatientCreateEmailOTPLimiter = RateLimiter{
	maxWrongAttemptsByIPperDay: 30, maxConsecutiveFailsByEmailAndIP: 5,
	client: config.PatientRedisClient, SlowBruteKeyPrefix: "patient_create_email_otp_fail_ip_per_day",
	ConsecutiveFailsKeyPrefix: "patient_create_email_otp_fail_consecutive_email_and_ip",
}

var PatientValidateEmailOTPLimiter = RateLimiter{
	maxWrongAttemptsByIPperDay: 20, maxConsecutiveFailsByEmailAndIP: 5,
	client: config.PatientRedisClient, SlowBruteKeyPrefix: "patient_validate_email_otp_fail_ip_per_day",
	ConsecutiveFailsKeyPrefix: "patient_validate_email_otp_fail_consecutive_email_and_ip",
}
