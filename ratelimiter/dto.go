package ratelimiter

import (
	"encoding/json"
	"time"
)

type failure struct {
	Attempts  int       `json:"attempts"`
	BlockTill time.Time `json:"block_until"`
}

func (t *failure) marshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *failure) unmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	return nil
}
