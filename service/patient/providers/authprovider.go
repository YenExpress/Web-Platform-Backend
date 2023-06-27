package providers

import (
	"YenExpress/config"
	"YenExpress/service/dto"
	"YenExpress/service/guard"
	"YenExpress/service/postoffice"

	"YenExpress/helper"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Joker666/cogman/util"
	limiter "github.com/codeNino/ratelimiter"

	"github.com/go-redis/redis"
)

func userExists(id uint) bool {
	var user *dto.Patient
	err := config.DB.Where("ID = ?", id).First(&user).Error
	if err != nil {
		return false
	}
	return true
}

type LoginManager struct {
}

func (manager *LoginManager) GenerateAuthTokens(user dto.Patient) (identifier, accessToken, refreshToken string) {
	var token1, token2, token3 string
	var wg sync.WaitGroup
	wg.Add(3)

	go func(identifier *string) {
		*identifier = JWTMaker.CreateIdentifier(
			guard.Identifier{
				Issuer: config.ServerDomain, Audience: config.WebClientDomain,
				Subject: "Token Bearing User Identity", UserId: user.ID,
				UserName: user.UserName, LastName: user.LastName,
				FirstName: user.FirstName, Email: user.Email, Sex: user.Sex,
			})
		wg.Done()
	}(&token1)

	go func(token *string) {
		*token = JWTMaker.CreateToken(
			guard.Bearer{
				UserId:     user.ID,
				Expiration: time.Now().Add(time.Hour * 24 * 3),
				Issuer:     config.ServerDomain, Class: "access_token",
			})
		wg.Done()
	}(&token2)

	go func(token *string) {
		*token = JWTMaker.CreateToken(
			guard.Bearer{
				UserId:     user.ID,
				Expiration: time.Now().Add(time.Hour * 24 * 30),
				Issuer:     config.ServerDomain, Class: "refresh_token",
			})
		wg.Done()
	}(&token3)
	wg.Wait()
	identifier, accessToken, refreshToken = token1, token2, token3
	return
}

type RegistrationManager struct {
	client            *redis.Client
	validationMailKey string
}

func (manager *RegistrationManager) AuthNewEmail(email string) {
	handler := signUpMailHandler
	handler.Mailer.To = email
	handler.Mailer.MailBodyVal = postoffice.OneTimePassword{
		OTP:       helper.GenerateOTPCode(7),
		Validity:  "10 minutes",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 10),
	}
	tokenized, _ := handler.Mailer.MailBodyVal.Marshal()
	manager.client.HSet(manager.validationMailKey, email, tokenized)
	task, err := handler.GetTask()
	if err != nil {
		log.Println(err)
		return
	}
	handlerfunc := util.HandlerFunc(handler.DoTask)
	if err != nil {
		log.Println(err)
		return
	}
	helper.QueueTask(task, handlerfunc)
}

func (manager *RegistrationManager) EnableSignUpwithMail(email, otp string) error {
	if manager.client.HExists(manager.validationMailKey, email).Val() {
		var auth_details postoffice.OneTimePassword
		auth_details.Unmarshal([]byte(manager.client.HGet(manager.validationMailKey, email).Val()))
		if auth_details.OTP == otp && time.Now().Before(auth_details.ExpiresAt) {
			manager.client.HDel(manager.validationMailKey, email)
			return nil
		}
		return errors.New("Invalid OTP")
	}
	return errors.New("Incorrect Email Provided")
}

var (
	LoginService LoginManager = LoginManager{}

	RegistrationService RegistrationManager = RegistrationManager{client: config.PatientRedisClient,
		validationMailKey: "NewPatientEmailValidation",
	}

	JWTMaker = &guard.JWTStrategy{
		SecretKey: config.JwtSecret,
		UserValid: userExists}

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
