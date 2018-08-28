/**
 *	Author Fanxu(746439274@qq.com)
 */

package queue

import (
	"asyncTask/model"
	"time"
	"asyncTask/config"
	"encoding/json"
	"log"
)

var Queue struct{
	Tasks chan model.Command
	Worker func(task model.Command) error
	Consumer func()
}
func Init() {
	Queue.Tasks = make(chan model.Command,1024)
	Queue.Worker = func(task model.Command) error {
		timestamp := time.Now().Unix()
		//格式化为字符串,tm为Time类型
		tm := time.Unix(timestamp, 0)
		task.Time = tm.Format("2006-01-02 03:04:05 PM")
		Queue.Tasks <- task
		return nil
	}
	Queue.Consumer = func() {
		for {
			select {
			case task := <- Queue.Tasks:
				//处理过来的队列请求
				switch task.Action {
				case config.ACTION_RECONNECT_MYSQL:
					jsonStr,err:= json.Marshal(task.Data)
					if err != nil {
						log.Println(err.Error())
						break
					}
					var d model.ReconnectData
					json.Unmarshal(jsonStr,&d)
					break
				case config.ACTION_PRINT:
					jsonStr,err:= json.Marshal(task.Data)
					if err != nil {
						log.Println(err.Error())
						break
					}
					log.Println(string(jsonStr))
					break
				}
			}
		}
	}
}