/**
 *	Author Fanxu(746439274@qq.com)
 */
package control

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"asyncTask/helpers"
	"strings"
	"fmt"
	"log"
	"asyncTask/config"
	"asyncTask/model"
)

type ProgramParams struct {
	Port string `json:"port"`
}

func (this * BaseControl) ProgramStart(w http.ResponseWriter, r *http.Request){
	p,err := getCurrentParam(r)
	if err != nil {
		panic(model.CustomException{402,err.Error()})
	}
	if p.Command == "" {
		panic(model.CustomException{402,"command not found"})
	}
	if p.Id=="" || !helpers.IsDirExists(config.All().UploadPath+p.Id) {
		panic(model.CustomException{402,"目录不存在"})
	}
	res,err2 := helpers.ExecShell("cd " + config.All().UploadPath+p.Id + " && " + p.Command)
	log.Println(res)
	if err2 != nil || strings.Contains(res,"command not found"){
		//同步返回
		panic(model.CustomException{400,"请先确保软件正确安装"+p.Command})
	}
	//同步返回
	helpers.ResponseJson(200,"ok!",map[string]interface{}{"logs":res},w)
	return
}

func (this * BaseControl) ProgramStatus(w http.ResponseWriter, r *http.Request){
	p,err := getParam(r)
	if err != nil {
		helpers.ResponseJson(402,err.Error(),helpers.EmptyResult,w)
		return
	}
	if p.Port == "" {
		panic(model.CustomException{400,"command not found"})
	}
	res,err2 := helpers.ExecShell(fmt.Sprintf("telnet 127.0.0.1 %d",p.Port))
	log.Println(res)
	if err2 != nil || strings.Contains(res,"command not found") || strings.Contains(res,"Connection refused"){
		//同步返回
		helpers.ResponseJson(200,"error",map[string]interface{}{"status":0},w)
		return
	}
	helpers.ResponseJson(200,"ok",map[string]interface{}{"status":1},w)
	return
}

func getParam(r *http.Request) (ProgramParams,error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		return ProgramParams{},err
	}
	var p ProgramParams
	err1 := json.Unmarshal(body,&p)
	if err1 != nil{
		return ProgramParams{},err
	}
	return p,nil
}
