package taskmaster

import (
	"YenExpress/config"

	"log"

	"github.com/Joker666/cogman"
	"github.com/Joker666/cogman/util"
)

func StartTaskMaster() {
	if err := cogman.StartBackground(config.TaskMasterCfg); err != nil {
		log.Fatal(err.Error())
	}

}

func QueueTask(task *util.Task, handler util.Handler) {
	if err := cogman.SendTask(*task, handler); err != nil {
		log.Fatal(err.Error())
	}

}

func SendNewAccountOTP(address string) {
	task, err := GetNewAccountOTPMailTask(address)
	if err != nil {
		return
	}
	handlerfunc := util.HandlerFunc(GetNewAccountOTPMailHandler)
	if err != nil {
		return
	}
	QueueTask(task, handlerfunc)
}
