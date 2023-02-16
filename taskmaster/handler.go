package taskmaster

import (
	"context"
	"encoding/json"
	"log"
	// "github.com/Joker666/cogman/util"
)

type MailTaskBody struct {
	address string
}

func GetNewAccountOTPMailHandler(ctx context.Context, payload []byte) error {

	var body MailTaskBody
	if err := json.Unmarshal(payload, &body); err != nil {
		log.Print("new account otp process error", err)
		return err
	}
	// otp := ""
	return nil
}
