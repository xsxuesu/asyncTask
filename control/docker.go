/**
 *	Author Fanxu(746439274@qq.com)
 */
package control

import (
	"net/http"
	"asyncTask/helpers"
	"strings"
	"io/ioutil"
	"encoding/json"
	"asyncTask/config"
	"fmt"
	"asyncTask/model"
)

type DockerPerformanceResult struct {
	Name string `json:"name"`
	ContainerId string `json:"container_id"`
	CpuUsage string `json:"cpuUsage"`
	MemUsage string `json:"memUsage"`
	Status int `json:"status"`
}

func (this * BaseControl) DockerPerformance(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		panic(model.CustomException{402, err.Error()})
	}
	fmt.Println(string(body))
	var p SystemPerformanceParam
	err = json.Unmarshal(body,&p)
	if err != nil{
		panic(model.CustomException{402, "dockers params : "+err.Error()})
	}
	dockers := make(map[string]*DockerPerformanceResult,0)
	contents , err := helpers.ExecShell("docker ps -a")
	if err != nil{
		panic(model.CustomException{400,err.Error()})
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if len(fields) < 1{
			continue
		}
		if "NAMES" == fields[len(fields)-1]{
			continue
		}
		fmt.Println(fields[len(fields)-1])
		if !helpers.HasElement(p.Dockers,fields[len(fields)-1]){
			continue
		}
		dockers[fields[0]] = &DockerPerformanceResult{
			fields[len(fields)-1],
			fields[0],
			"0.00%",
			"0.00%",
			0,
		}
	}
	contents , err = helpers.ExecShell("docker stats --no-stream")
	if err != nil{
		panic(model.CustomException{400,err.Error()})
	}
	lines = strings.Split(string(contents), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if len(fields) < 1{
			continue
		}
		if "PIDS" == fields[len(fields)-1]{
			continue
		}
		_, ok := dockers[fields[0]]
		if ok {
			dockers[fields[0]].Status = 1
			dockers[fields[0]].CpuUsage = fields[1]
			dockers[fields[0]].MemUsage = fields[7]
		}
	}
	helpers.ResponseJson(200,"ok",map[string] interface{}{
		"list":dockers,
	},w)
	return
}

type DockerExecuteParam struct {
	Type string `json:"type"`
	ContainerId string `json:"container_id"`
}
func (this * BaseControl) DockerExecute(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	if err != nil{
		panic(model.CustomException{402,"params error!"})
	}
	var u DockerExecuteParam
	err = json.Unmarshal(body,&u)
	fmt.Println(u)
	if err != nil{
		panic(model.CustomException{402,"params error!"})
	}
	var cmd string
	switch u.Type {
	case config.DOCKER_EXEC_START:
		cmd = fmt.Sprintf("docker start %s",u.ContainerId)
			break
	case config.DOCKER_EXEC_STOP:
		cmd = fmt.Sprintf("docker stop %s",u.ContainerId)
			break
	case config.DOCKER_EXEC_RESTART:
		cmd = fmt.Sprintf("docker restart %s",u.ContainerId)
			break
	case config.DOCKER_EXEC_BACKUP:
		_,err = helpers.ExecShell("docker stop "+u.ContainerId)
		if err != nil{
			panic(model.CustomException{411,"execute error!"})
		}
		cmd = fmt.Sprintf("docker export --output='%s' %s",u.ContainerId,u.ContainerId)
			break
	default:
		panic(model.CustomException{402,"type error!"})
		break
	}
	_,err = helpers.ExecShell(cmd)
	if err != nil{
		panic(model.CustomException{410,"execute error!"})
	}
	switch u.Type {
	case config.DOCKER_EXEC_BACKUP:
		_,err = helpers.ExecShell("docker start "+u.ContainerId)
		if err != nil{
			panic(model.CustomException{411,"execute error!"})
		}
		break
	}
	helpers.ResponseJson(200,"ok",helpers.EmptyResult,w)
	return
}
