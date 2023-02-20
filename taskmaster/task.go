package taskmaster

import (
	"encoding/json"
	"log"

	"github.com/Joker666/cogman/util"
)

func parseBody(body interface{}) ([]byte, error) {
	pld, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return pld, nil
}

func GetNewAccountOTPMailTask(email string) (*util.Task, error) {
	body := MailTaskBody{
		address: email,
	}

	pld, err := parseBody(body)
	if err != nil {
		log.Print("Parse: ", err)
		return nil, err
	}

	task := &util.Task{
		Name:     "Send_new_account_otp_to_email",
		Payload:  pld,
		Priority: util.TaskPriorityHigh,
		Retry:    5,
	}

	return task, nil
}

func GetRateLimitingTask(ip_addr, email string) (*util.Task, error) {
	body := RateLimitingTaskBody{
		email_address: email,
		ip_address:    ip_addr,
	}

	pld, err := parseBody(body)
	if err != nil {
		log.Print("Parse: ", err)
		return nil, err
	}

	task := &util.Task{
		Name:     "",
		Payload:  pld,
		Priority: util.TaskPriorityHigh,
		Retry:    5,
	}

	return task, nil
}
