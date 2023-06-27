package models

import (
	"encoding/json"
	"time"
)

type CreateDTO struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" validate:"required,Enum=admin_superadmin"`
}

type SessionData struct {
	SessionID string    `json:"sessionID"`
	IPAddress string    `json:"ipAddr,omitempty"`
	Email     string    `json:"Email"`
	LoggedIn  time.Time `json:"loggedIn"`
	UserID    uint      `json:"userID"`
}

type SessionStore struct {
	Session map[string]SessionData
}

func (s *SessionStore) Encode() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SessionStore) Decode(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return nil
}
