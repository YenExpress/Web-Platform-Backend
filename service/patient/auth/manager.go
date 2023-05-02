package auth

import (
	// "context"
	"strconv"

	"YenExpress/config"
	"YenExpress/postoffice"
	"YenExpress/service/patient/guard"
	"YenExpress/toolbox"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Joker666/cogman/util"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

var loginMailHandler = postoffice.MailHandler[postoffice.OneTimePassword]{
	Mailer: postoffice.PostMan[postoffice.OneTimePassword]{
		MailTemplatePath: "/postoffice/templates/concurrentLoginValidation.html",
		Subject:          "Login Validation",
	},
}

var signUpMailHandler = postoffice.MailHandler[postoffice.OneTimePassword]{
	Mailer: postoffice.PostMan[postoffice.OneTimePassword]{
		MailTemplatePath: "/postoffice/templates/newAccountEmailValidation.html",
		Subject:          "New Account Email Validation",
	},
}

type loginManager struct {
	sessionKey        string
	client            *redis.Client
	validationMailKey string
}

// cache user session data upon successful login and authentication
// return error and sessionID value
func (manager *loginManager) createNewSession(user Patient, ip_addr string) string {
	sessID := uuid.New().String()
	data := sessionStore{
		Session: map[string]sessionData{
			sessID: sessionData{
				SessionID: sessID, UserID: user.ID,
				LoggedIn: time.Now(), Email: user.Email,
				IPAddress: ip_addr,
			},
		},
	}
	tokenized, _ := data.encode()
	manager.client.HSet(manager.sessionKey, strconv.FormatUint(uint64(user.ID), 10), tokenized).Val()
	return sessID
}

// add another session initiated by user on a different device
// return error and sessionID value
func (manager *loginManager) addSession(user Patient, ip_addr string) string {
	userKey := strconv.FormatUint(uint64(user.ID), 10)
	if manager.client.HExists(manager.sessionKey, userKey).Val() {
		var cookie sessionStore
		cookie.decode([]byte(manager.client.HGet(manager.sessionKey, userKey).Val()))
		sessID := uuid.New().String()
		cookie.Session[sessID] = sessionData{
			SessionID: sessID, IPAddress: ip_addr,
			Email: user.Email, LoggedIn: time.Now(), UserID: user.ID,
		}
		tokenized, _ := cookie.encode()
		manager.client.HSet(manager.sessionKey, userKey, tokenized).Val()
		return sessID
	}
	return manager.createNewSession(user, ip_addr)
}

// check if user is logged in and has an active session
func (manager *loginManager) checkActiveSession(userID uint) bool {
	return manager.client.HExists(manager.sessionKey, strconv.FormatUint(uint64(userID), 10)).Val()
}

func (manager *loginManager) endSession(userID uint, sessionID string) sessionData {
	userKey := strconv.FormatUint(uint64(userID), 10)
	if manager.client.HExists(manager.sessionKey, userKey).Val() {
		var cookie sessionStore
		cookie.decode([]byte(manager.client.HGet(manager.sessionKey, userKey).Val()))
		session_data := cookie.Session[sessionID]
		delete(cookie.Session, sessionID)
		if len(cookie.Session) == 0 {
			manager.client.HDel(manager.sessionKey, userKey)
			return session_data
		}
		tokenized, _ := cookie.encode()
		manager.client.HSet(manager.sessionKey, userKey, tokenized).Val()
		return session_data
	}
	return sessionData{}
}

func (manager *loginManager) authConcurrentSignin(email string) {
	handler := loginMailHandler
	handler.Mailer.To = email
	handler.Mailer.MailBodyVal = postoffice.OneTimePassword{
		OTP:       toolbox.GenerateOTPCode(7),
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
	toolbox.QueueTask(task, handlerfunc)
}

func (manager *loginManager) enableConcurrentSignin(user Patient, ip_addr, otp string) (identifier, accessToken, refreshToken string, err error) {
	if manager.client.HExists(manager.validationMailKey, user.Email).Val() {
		var auth_details postoffice.OneTimePassword
		auth_details.Unmarshal([]byte(manager.client.HGet(manager.sessionKey, user.Email).Val()))
		if auth_details.OTP == otp && time.Now().Before(auth_details.ExpiresAt) {
			sessID := manager.addSession(user, ip_addr)
			identifier, accessToken, refreshToken := manager.generateAuthTokens(user, sessID)
			manager.client.HDel(manager.validationMailKey, user.Email)
			return identifier, accessToken, refreshToken, nil

		}
		return "", "", "", errors.New("Invalid OTP")

	}
	return "", "", "", errors.New("Incorrect Email")

}

func (manager *loginManager) generateAuthTokens(user Patient, sessionID string) (identifier, accessToken, refreshToken string) {
	var token1, token2, token3 string
	var wg sync.WaitGroup
	wg.Add(3)

	go func(identifier *string) {
		*identifier = guard.Bearer.CreateIdentifier(
			guard.Identifier{
				Issuer: config.ServerDomain, Audience: config.WebClientDomain,
				Subject: "Token Bearing User Identity", UserId: user.ID,
				UserName: user.UserName, LastName: user.LastName,
				FirstName: user.FirstName, Email: user.Email, Sex: user.Sex,
			})
		wg.Done()
	}(&token1)

	go func(token *string) {
		*token = guard.Bearer.CreateToken(
			guard.Payload{
				UserId: user.ID, SessionID: sessionID,
				Expiration: time.Now().Add(time.Hour * 24 * 3),
				Issuer:     config.ServerDomain, Class: "access_token",
			})
		wg.Done()
	}(&token2)

	go func(token *string) {
		*token = guard.Bearer.CreateToken(
			guard.Payload{
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

type signUpManager struct {
	client            *redis.Client
	validationMailKey string
}

func (manager *signUpManager) authNewEmail(email string) {
	handler := signUpMailHandler
	handler.Mailer.To = email
	handler.Mailer.MailBodyVal = postoffice.OneTimePassword{
		OTP:       toolbox.GenerateOTPCode(7),
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
	toolbox.QueueTask(task, handlerfunc)
}

func (manager *signUpManager) enableSignUpwithMail(email, otp string) error {
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
	PatientLoginManager loginManager = loginManager{sessionKey: "ActivePatientSession",
		client: config.PatientRedisClient, validationMailKey: "concurrentPatientSessionValidation",
	}

	PatientRegistrationManager signUpManager = signUpManager{client: config.PatientRedisClient,
		validationMailKey: "NewPatientEmailValidation",
	}
)
