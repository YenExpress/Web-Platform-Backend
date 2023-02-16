package config

// import (
// 	// "YenExpress/models"
// 	// "context"
// 	"encoding/json"
// 	"time"

// 	"github.com/go-redis/redis"
// 	// "github.com/go-redis/redis"
// )

// type sessiondata struct {
// 	FirstName string    `json:"FirstName"`
// 	LastName  string    `json:"LastName"`
// 	Sex       string    `json:"Sex"`
// 	Email     string    `json:"Email"`
// 	TokenID   string    `json:"TokenID"`
// 	LoggedIn  time.Time `json:"LoggedIn"`
// }

// type SessionManager struct {
// 	session_key string
// 	client      *redis.Client
// }

// func (t *sessiondata) marshalBinary() ([]byte, error) {
// 	return json.Marshal(t)
// }

// func (t *sessiondata) unmarshalBinary(data []byte) error {
// 	if err := json.Unmarshal(data, &t); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // cache user session data upon successful login and authentication
// func (manager *SessionManager) RegisterSession(user models.ReadUser, token_id string) bool {
// 	if err := manager.client.Ping().Err(); err != nil {
// 		return false
// 	}
// 	data := sessiondata{
// 		Email: user.Email, TokenID: token_id,
// 		FirstName: user.FirstName, LastName: user.LastName,
// 		LoggedIn: user.LastLogin, Sex: user.Sex,
// 	}
// 	tokenized, _ := data.marshalBinary()
// 	result := manager.client.HSet(manager.session_key, user.ID, tokenized).Val()
// 	return result
// }

// // update cached user session data
// func (manager *SessionManager) UpdateSession(user_id string, update map[string]interface{}) bool {
// 	if manager.client.HExists(manager.session_key, user_id).Val() {
// 		var datastruct sessiondata
// 		err := datastruct.unmarshalBinary([]byte(manager.client.HGet(manager.session_key, user_id).Val()))
// 		if err != nil {
// 			return false
// 		}
// 		for session_key, session_value := range update {
// 			switch session_key {
// 			case "FirstName":
// 				datastruct.FirstName = session_value.(string)
// 			case "LastName":
// 				datastruct.LastName = session_value.(string)
// 			case "Sex":
// 				datastruct.Sex = session_value.(string)
// 			case "Email":
// 				datastruct.Email = session_value.(string)
// 			case "TokenID":
// 				datastruct.TokenID = session_value.(string)
// 			}
// 		}
// 		tokenized, _ := datastruct.marshalBinary()
// 		result := manager.client.HSet(manager.session_key, user_id, tokenized).Val()
// 		return result
// 	}
// 	return false
// }

// func (manager *SessionManager) GetSessionData(user_id string) (sessiondata, error) {
// 	if manager.client.HExists(manager.session_key, user_id).Val() {
// 		var datastruct sessiondata
// 		err := datastruct.unmarshalBinary([]byte(manager.client.HGet(manager.session_key, user_id).Val()))
// 		if err != nil {
// 			return sessiondata{}, err
// 		}
// 		return datastruct, nil
// 	}
// 	return sessiondata{}, nil
// }

// func (manager *SessionManager) EndSession(user_id string) (sessiondata, error) {
// 	if manager.client.HExists(manager.session_key, user_id).Val() {
// 		var datastruct sessiondata
// 		err := datastruct.unmarshalBinary([]byte(manager.client.HGet(manager.session_key, user_id).Val()))
// 		if err != nil {
// 			return datastruct, err
// 		}
// 		manager.client.HDel(manager.session_key, user_id)
// 		return datastruct, nil
// 	}
// 	return sessiondata{}, nil
// }

// var UserSessionManager SessionManager = SessionManager{session_key: "DAVS Users", client: RedisClient}
