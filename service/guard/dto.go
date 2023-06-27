package guard

import (
	"YenExpress/config"
	"time"
)

type Identifier struct {
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	Audience  string `json:"audience"`
	UserId    uint   `json:"userId"`
	UserName  string `json:"userName,omitempty"`
	Role      string `json:"role,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Sex       string `json:"sex"`
	Photo     string `json:"photo,omitempty"`
}

func (data *Identifier) Valid() error {
	if data.Issuer != config.ServerDomain || data.Audience != config.WebClientDomain {
		return ErrInvalidToken
	}
	return nil
}

type Bearer struct {
	UserId     uint      `json:"userId"`
	SessionID  string    `json:"sessionID"`
	Expiration time.Time `json:"expiration"`
	Issuer     string    `json:"issuer"`
	Class      string    `json:"class"`
}

func (payload *Bearer) Valid() error {
	if time.Now().After(payload.Expiration) {
		return ErrExpiredToken
	} else if payload.Issuer != config.ServerDomain {
		return ErrInvalidToken
	}
	return nil
}
