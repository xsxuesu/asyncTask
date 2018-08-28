/**
 *	Author Fanxu(746439274@qq.com)
 */
package control

import (
	"strconv"
	"strings"
	"log"
	"net/http"
	"path"
	"fmt"
	"os"
	"io"
	"asyncTask/helpers"
	"io/ioutil"
	"encoding/json"
	"asyncTask/model"
)

type ChaincodeUploadByGitParam struct {
	ChaincodeName string `json:"chaincodeName"`
	ChaincodePath string `json:"chaincodePath"`
	ChaincodeGit string `json:"chaincodeGit"`
} 
func (this * BaseControl) ChaincodeUploadByZip(w http.ResponseWriter, r *http.Request){
	if "POST" != r.Method {
		showUploadHtml(w)
		return
	}
	file, header, err := r.FormFile("chaincodeZip")
	if err != nil{
		helpers.ResponseJson(402,err.Error(),helpers.EmptyResult,w)
		return
	}
	fileSuffix := path.Ext(header.Filename) //获取文件后缀
	fmt.Println(fileSuffix)
	if fileSuffix != ".zip" {
		panic(model.CustomException{406,"目前只允许zip文件!"})
	}
	fpath := r.FormValue("chaincodePath")
	if fpath == "" {
		panic(model.CustomException{402,"path 不能为空"})
	}
	name := r.FormValue("chaincodeName")
	if name == "" {
		panic(model.CustomException{402,"name 不能为空"})
	}
	name = ""
	defer file.Close()
	if !helpers.IsDirExists(fpath+"/"+name) {
		err := os.MkdirAll(fpath+"/"+name,os.ModePerm)
		if err != nil {
			panic(model.CustomException{403,err.Error()})
		}
	}
	f,err:=os.Create("/tmp/"+header.Filename)
	if err != nil{
		panic(model.CustomException{403,err.Error()})
	}
	defer f.Close()
	_,err = io.Copy(f,file)
	if err != nil {
		panic(model.CustomException{403,err.Error()})
	}
	//_,err2 := helpers.ExecShell("cd /tmp/ && unzip -o "+header.Filename + " -d "+fpath+"/"+name+" && rm -rf "+header.Filename)
	tmpPath := "/tmp/"+r.FormValue("chaincodeName")
	fmt.Println(tmpPath)
	err = os.MkdirAll(tmpPath,os.ModePerm)
	if err != nil{
		panic(model.CustomException{405,err.Error()})
	}
	_,err = helpers.ExecShell("cd /tmp/ && unzip -o "+header.Filename + " -d "+tmpPath+" && rm -rf "+header.Filename)
	if err != nil {
		panic(model.CustomException{405,"解压文件失败:"+err.Error()})
	}
	result,err3 := helpers.ExecShell("ls " + tmpPath)
	if err3 != nil {
		panic(model.CustomException{405,"文件错误:"+err3.Error()})
		return
	}
	if strings.Contains(result,"main.go"){
		panic(model.CustomException{405,"智能合约格式不正确:"})
	}
	//智能合约里的文件夹数量
	cmd := fmt.Sprintf("ls %s -l |grep '^d'|wc -l",tmpPath+"/"+name)
	s,_ := helpers.ExecShell(cmd)
	dirCount,_:= strconv.Atoi(strings.Replace(s,"\n","",1))
	//智能合约里的文件数量
	cmd = fmt.Sprintf("ls %s -l |grep '^-'|wc -l",tmpPath+"/"+name)
	s,_ = helpers.ExecShell(cmd)
	fileCount,_:= strconv.Atoi(strings.Replace(s,"\n","",1))
	newDirName := ""
	if fileCount<=0 && dirCount == 1{
		newDirName = strings.Replace(result,"\n","",1)
		if !helpers.IsDirExists(fpath+"/"+newDirName){
			os.MkdirAll(fpath+"/"+newDirName,os.ModePerm)
		}else{
			helpers.ExecShell("rm -rf " + fpath+"/"+newDirName+"/*")
		}
		//将文件夹子目录移到当前目录
		_,err = helpers.ExecShell("cp -rf " + tmpPath+"/"+newDirName+"/* "+fpath+"/"+newDirName+" && rm -rf "+tmpPath)
		if err != nil {
			panic(model.CustomException{405,"解压文件失败:"+err.Error()})
		}
		result,_= helpers.ExecShell("ls " + fpath+"/"+newDirName)
	}
	if !strings.Contains(result,"main.go"){
		panic(model.CustomException{405,"智能合约格式不正确"})
	}
	//同步返回
	helpers.ResponseJson(200,"上传文件成功!",map[string] interface{}{ "path":fpath+"/"+newDirName},w)
	return
}

func (this * BaseControl) ChaincodeUploadByGit(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	log.Println(string(body))
	if err != nil{
		panic(model.CustomException{402,"params error"})
	}
	var u ChaincodeUploadByGitParam
	err = json.Unmarshal(body,&u)
	if err != nil{
		panic(model.CustomException{402,"params error"+err.Error()})
	}
	//cpath := u.ChaincodePath+"/"+u.ChaincodeName
	cpath := u.ChaincodePath
	tmpPath := "/tmp/"+u.ChaincodeName
	os.MkdirAll(tmpPath,os.ModePerm)
	command := "cd "+tmpPath+" && git clone "+u.ChaincodeGit
	res,err :=helpers.ExecShell(command)
	if err != nil{
		helpers.ResponseJson(500,err.Error(),helpers.EmptyResult,w)
		return
	}
	result,err3 := helpers.ExecShell("ls " + tmpPath)
	if err3 != nil {
		panic(model.CustomException{500,"文件错误"+err3.Error()})
	}
	newDirName := strings.Replace(result,"\n","",1)
	if !helpers.IsDirExists(cpath+"/"+newDirName){
		os.MkdirAll(cpath+"/"+newDirName,os.ModePerm)
	}else{
		helpers.ExecShell("rm -rf " + cpath+"/"+newDirName+"/*")
	}
	_,err = helpers.ExecShell("cp -rf " + tmpPath+"/"+newDirName+"/* "+cpath+"/"+newDirName+" && rm -rf "+tmpPath)
	if err != nil {
		panic(model.CustomException{500,"上传失败"+err.Error()})
	}
	result,_= helpers.ExecShell("ls " + cpath+"/"+newDirName)
	if !strings.Contains(result,"main.go"){
		panic(model.CustomException{500,"智能合约格式不正确"})
	}
	helpers.ResponseJson(200,"ok",map[string] interface{}{"logs":res,"path":cpath+newDirName},w)
}