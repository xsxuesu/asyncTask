/**
 *	Author Fanxu(746439274@qq.com)
 */
package control

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"asyncTask/helpers"
	"asyncTask/config"
	"asyncTask/model"
)

type SdkInstallParam struct {
	SdkPath string `json:"sdkPath"`
	Id string `json:"id"`
	SdkGit string `json:"sdkGit"`
	SdkConfig []byte `json:"sdkConfig"`
}

type SdkExecuteParam struct {
	ShellContent []byte `json:"shellContent"`
}

func (this * BaseControl) SdkInstall(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil{
		panic(model.CustomException{402,"params error"})
	}
	var u SdkInstallParam
	err1 := json.Unmarshal(body,&u)
	if err1 != nil{
		panic(model.CustomException{402,"params error"})
	}
	var command = "cd "+u.SdkPath+" && git pull origin master && rm -rf config.json && echo '"+ string(u.SdkConfig)+"' >>config.json && pm2 kill && pm2 start app"
	if !helpers.IsDirExists(u.SdkPath){
		err := os.MkdirAll(config.All().UploadPath, os.ModePerm)
		if err != nil{
			panic(model.CustomException{500,u.SdkPath + err.Error()})
		}
		command = "git clone "+u.SdkGit + " "+u.SdkPath+ " && cd "+u.SdkPath+" && echo '"+ string(u.SdkConfig)+"' >>config.json && npm install  && pm2 kill && pm2 start app"
	}
	res,err :=helpers.ExecShell(command)
	if err != nil{
		panic(model.CustomException{500, err.Error()})
	}
	helpers.ResponseJson(200,"ok",map[string] interface{}{"logs":res},w)
}
