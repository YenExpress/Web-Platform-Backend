package toolbox

import (
	"YenExpress/config"

	"fmt"
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

func NameTask(task string, user string) string {
	return fmt.Sprintf(
		"%v_for_%v_taggedBy_%v", task, user, GenerateRandStr(10))
}
