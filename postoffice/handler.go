package postoffice

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Joker666/cogman/util"
	"github.com/google/uuid"
)

type MailHandler[T any] struct {
	Mailer PostMan[T]
}

func (h *MailHandler[T]) encodeHandler() ([]byte, error) {
	return json.Marshal(h)
}

func (h *MailHandler[T]) DecodeHandler(data []byte) error {
	if err := json.Unmarshal(data, &h); err != nil {
		return err
	}
	return nil
}

func (h *MailHandler[T]) GetTask() (*util.Task, error) {
	pld, err := h.encodeHandler()
	if err != nil {
		log.Print("Parse: ", err)
		return nil, err
	}

	task := &util.Task{
		Name:     uuid.New().String(),
		Payload:  pld,
		Priority: util.TaskPriorityHigh,
		Retry:    5,
	}

	return task, nil
}

func (h *MailHandler[T]) DoTask(ctx context.Context, payload []byte) error {

	if err := h.DecodeHandler(payload); err != nil {
		log.Print("process error ==> ", err)
		return err
	}

	resp, err := h.Mailer.SendMail()
	if err != nil {
		log.Print("error ==> ", err)
		return err
	} else if resp.StatusCode != 202 {
		log.Printf("response from sender api ==> %v with statusCode %v", resp.Body, resp.StatusCode)
		return errors.New(resp.Body)
	}
	return nil
}
