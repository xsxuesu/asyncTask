/**
 *	Author Fanxu(746439274@qq.com)
 */
package control

import (
	"asyncTask/helpers"
	"asyncTask/model"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type EnvironmentParams struct {
	Command string `json:"command"`
	Id      string `json:"id"`
}

func (this *BaseControl) EnvironmentIp(w http.ResponseWriter, r *http.Request) {
	res, err := helpers.ExecShell("docker ps -a")
	log.Println(res)
	if err != nil {
		helpers.ResponseJson(407, "请查看日志错误信息:", map[string]interface{}{"logs": res}, w)
		return
	}
	if strings.Contains(res, "hyperledger") {
		helpers.ResponseJson(407, "此服务器已经部署了其他项目，请查看其他服务器！", map[string]interface{}{"logs": res}, w)
		return
	}
	//同步返回
	helpers.ResponseJson(200, "环境安装完成!", map[string]interface{}{"logs": res}, w)
	return
}

func (this *BaseControl) EnvironmentCheck(w http.ResponseWriter, r *http.Request) {
	p, err := getCurrentParam(r)
	if err != nil {
		panic(model.CustomException{402, err.Error()})
	}
	if p.Command == "" {
		panic(model.CustomException{402, "command not found "})
	}
	res, err2 := helpers.ExecShell(p.Command)
	log.Println(res)
	if err2 != nil || strings.Contains(res, "command not found") {
		//同步返回
		panic(model.CustomException{407, "software not install "})
	}
	//同步返回
	helpers.ResponseJson(200, "环境检测完成!", helpers.EmptyResult, w)
	return
}
func (this *BaseControl) EnvironmentInstall(w http.ResponseWriter, r *http.Request) {
	p, err := getCurrentParam(r)
	if err != nil {
		panic(model.CustomException{402, err.Error()})
	}
	if p.Command == "" {
		panic(model.CustomException{402, "command not found "})
	}
	res, err2 := helpers.ExecShell(p.Command)
	log.Println(res)
	if err2 != nil || strings.Contains(res, "command not found ") {
		panic(model.CustomException{402, "command not found "})
		//同步返回
		helpers.ResponseJson(407, "请查看日志错误信息 "+p.Command, map[string]interface{}{"logs": res}, w)
		return
	}
	//同步返回
	helpers.ResponseJson(200, "环境安装完成!", map[string]interface{}{"logs": res}, w)
	return
}

func getCurrentParam(r *http.Request) (EnvironmentParams, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return EnvironmentParams{}, err
	}
	var p EnvironmentParams
	err1 := json.Unmarshal(body, &p)
	if err1 != nil {
		return EnvironmentParams{}, err
	}
	return p, nil
}
