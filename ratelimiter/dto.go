package ratelimiter

import (
	"encoding/json"
	"time"
)

type LoginLimit struct {
	FailedAttempts int       `json:"failed_attempts"`
	BlockDuration  time.Time `json:"block_duration"`
}

func (t *LoginLimit) marshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *LoginLimit) unmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	return nil
}
