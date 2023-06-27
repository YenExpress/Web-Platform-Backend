package providers

import (
	"YenExpress/config"
	"YenExpress/helper"
	model "YenExpress/service/admin/models"
	"YenExpress/service/dto"
	"YenExpress/service/guard"
	"YenExpress/service/postoffice"
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Joker666/cogman/util"
	limiter "github.com/codeNino/ratelimiter"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func userExists(id uint) bool {
	var user *dto.Admin
	err := config.DB.Where("ID = ?", id).First(&user).Error
	if err != nil {
		return false
	}
	return true
}

type LoginManager struct {
	sessionKey        string
	client            *redis.Client
	validationMailKey string
}

// cache user session data upon successful login and authentication
// return error and sessionID value
func (manager *LoginManager) CreateNewSession(user dto.Admin, ip_addr string) string {
	id_str := strconv.FormatUint(uint64(user.ID), 10)
	sessID := uuid.New().String()
	data := model.SessionStore{
		Session: map[string]model.SessionData{
			sessID: model.SessionData{
				SessionID: sessID, UserID: user.ID,
				LoggedIn: time.Now(), Email: user.Email,
				IPAddress: ip_addr,
			},
		},
	}
	tokenized, _ := data.Encode()
	manager.client.HSet(manager.sessionKey, id_str, tokenized).Val()
	return sessID
}

// add another session initiated by user on a different device
// return error and sessionID value
func (manager *LoginManager) AddSession(user dto.Admin, ip_addr string) string {
	id_str := strconv.FormatUint(uint64(user.ID), 10)
	if manager.client.HExists(manager.sessionKey, id_str).Val() {
		var cookie model.SessionStore
		cookie.Decode([]byte(manager.client.HGet(manager.sessionKey, id_str).Val()))
		sessID := uuid.New().String()
		cookie.Session[sessID] = model.SessionData{
			SessionID: sessID, IPAddress: ip_addr,
			Email: user.Email, LoggedIn: time.Now(), UserID: user.ID,
		}
		tokenized, _ := cookie.Encode()
		manager.client.HSet(manager.sessionKey, id_str, tokenized).Val()
		return sessID
	}
	return manager.CreateNewSession(user, ip_addr)
}

// check if user is logged in and has an active session
func (manager *LoginManager) CheckActiveSession(userID uint) bool {
	id_str := strconv.FormatUint(uint64(userID), 10)
	return manager.client.HExists(manager.sessionKey, id_str).Val()
}

func (manager *LoginManager) EndSession(userID uint, sessionID string) model.SessionData {
	id_str := strconv.FormatUint(uint64(userID), 10)
	if manager.client.HExists(manager.sessionKey, id_str).Val() {
		var cookie model.SessionStore
		cookie.Decode([]byte(manager.client.HGet(manager.sessionKey, id_str).Val()))
		session_data := cookie.Session[sessionID]
		delete(cookie.Session, sessionID)
		if len(cookie.Session) == 0 {
			manager.client.HDel(manager.sessionKey, id_str)
			return session_data
		}
		tokenized, _ := cookie.Encode()
		manager.client.HSet(manager.sessionKey, id_str, tokenized).Val()
		return session_data
	}
	return model.SessionData{}
}

func (manager *LoginManager) AuthConcurrentSignin(email string) {
	handler := loginMailHandler
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

func (manager *LoginManager) EnableConcurrentSignin(user dto.Admin, ip_addr, otp string) (identifier, accessToken, refreshToken string, err error) {
	if manager.client.HExists(manager.validationMailKey, user.Email).Val() {
		var auth_details postoffice.OneTimePassword
		auth_details.Unmarshal([]byte(manager.client.HGet(manager.validationMailKey, user.Email).Val()))
		if auth_details.OTP == otp && time.Now().Before(auth_details.ExpiresAt) {
			sessID := manager.AddSession(user, ip_addr)
			identifier, accessToken, refreshToken := manager.GenerateAuthTokens(user, sessID)
			manager.client.HDel(manager.validationMailKey, user.Email)
			return identifier, accessToken, refreshToken, nil

		}
		return "", "", "", errors.New("Invalid OTP")

	}
	return "", "", "", errors.New("Incorrect Email")

}

func (manager *LoginManager) GenerateAuthTokens(user dto.Admin, sessionID string) (identifier, accessToken, refreshToken string) {
	var token1, token2, token3 string
	var wg sync.WaitGroup
	wg.Add(3)

	go func(identifier *string) {
		*identifier = JWTMaker.CreateIdentifier(
			guard.Identifier{
				Issuer: config.ServerDomain, Audience: config.WebClientDomain,
				Subject: "Token Bearing Admin Identity", UserId: user.ID,
				Role: user.Role, LastName: user.LastName,
				FirstName: user.FirstName, Email: user.Email,
			})
		wg.Done()
	}(&token1)

	go func(token *string) {
		*token = JWTMaker.CreateToken(
			guard.Bearer{
				UserId: user.ID, SessionID: sessionID,
				Expiration: time.Now().Add(time.Hour * 24 * 3),
				Issuer:     config.ServerDomain, Class: "access_token",
			})
		wg.Done()
	}(&token2)

	go func(token *string) {
		*token = JWTMaker.CreateToken(
			guard.Bearer{
				UserId: user.ID, SessionID: sessionID,
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
	LoginService LoginManager = LoginManager{sessionKey: "ActiveAdminSession",
		client: config.PatientRedisClient, validationMailKey: "ConcurrentAdminSessionValidation",
	}

	RegistrationService RegistrationManager = RegistrationManager{client: config.PatientRedisClient,
		validationMailKey: "NewAdminEmailValidation",
	}

	JWTMaker = &guard.JWTStrategy{
		SecretKey: config.JwtSecret,
		UserValid: userExists}

	LoginLimiter = limiter.RateLimiter{
		TotalLimit: 100, BurstLimit: 10, MaxTime: time.Hour * 24, BurstPeriod: time.Hour * 1,
		Client: config.PatientRedisClient, TotalLimitPrefix: "admin_login_fail_ip_per_day",
		BurstLimitPrefix: "admin_login_fail_consecutive_email_and_ip",
	}

	CreateOTPLimiter = limiter.RateLimiter{
		TotalLimit: 30, BurstLimit: 5, MaxTime: time.Hour * 24, BurstPeriod: time.Minute * 30,
		Client: config.PatientRedisClient, TotalLimitPrefix: "admin_create_otp_fail_per_day",
		BurstLimitPrefix: "admin_create_otp_fail_consecutive",
	}

	EmailValidationLimiter = limiter.RateLimiter{
		TotalLimit: 30, BurstLimit: 5, MaxTime: time.Hour * 24, BurstPeriod: time.Minute * 30,
		Client: config.PatientRedisClient, TotalLimitPrefix: "admin_validate_email_otp_fail_per_day",
		BurstLimitPrefix: "admin_validate_email_otp_fail_consecutive",
	}
)
