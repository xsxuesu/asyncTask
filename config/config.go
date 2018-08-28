/**
 *	Author Fanxu(746439274@qq.com)
 */

package config

import (
	"asyncTask/model"
	"log"
	"io/ioutil"
	"encoding/json"
)

const (
	ACTION_PRINT = 1
	ACTION_RECONNECT_MYSQL = 2
	DOCKER_EXEC_START = "0"
	DOCKER_EXEC_STOP = "1"
	DOCKER_EXEC_RESTART = "2"
	DOCKER_EXEC_BACKUP = "3"
)

var __config * model.Config

func All() * model.Config{
	return __config
}

func InitConf(args []string) {
	var confFile = "C:/work/fanxu/go/src/asyncTask/config/online.json"
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-f":
			if i == len(args) {
				log.Fatalln("invalid config file")
			}
			confFile = args[i+1]
		}
	}
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatalln(err.Error(), "cannot find the file: online.json")
	}
	err = json.Unmarshal(data, &__config)
	if err != nil {
		log.Fatalln(err.Error(), "cannot parse the file: online.json")
	}
}

