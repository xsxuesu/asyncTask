/**
 *	Author Fanxu(746439274@qq.com)
 */
package control

import (
	"net/http"
	"asyncTask/helpers"
	"os"
	"io"
	"asyncTask/config"
	"asyncTask/model"
)
// 获取文件大小
type Sizer interface {
	Size() int64
}
func (this * BaseControl) FileUpload(w http.ResponseWriter, r *http.Request){
	if "POST" != r.Method {
		showUploadHtml(w)
		return
	}
	file, header, err := r.FormFile("certTar")
	if err != nil{
		panic(model.CustomException{402,err.Error()})
	}
	id := r.FormValue("id")
	if id == "" {
		panic(model.CustomException{402,"id 不能为空"})
	}
	defer file.Close()
	if !helpers.IsDirExists(config.All().UploadPath) {
		err := os.Mkdir(config.All().UploadPath, os.ModePerm)
		if err != nil{
			panic(model.CustomException{403,err.Error()})
		}
	}
	f,err:=os.Create(config.All().UploadPath+"/"+header.Filename)
	if err != nil{
		panic(model.CustomException{500,err.Error()})
	}
	defer f.Close()
	_,err = io.Copy(f,file)
	if err != nil {
		panic(model.CustomException{500,err.Error()})
	}
	_,err = helpers.ExecShell("cd "+config.All().UploadPath+" && rm -rf "+id+" && tar -xvf "+header.Filename)
	if err != nil {
		panic(model.CustomException{500,err.Error()})
	}
	//同步返回
	helpers.ResponseJson(200,"上传文件并解压成功!",helpers.EmptyResult,w)
	return
}
func (this * BaseControl) FileUploadById(w http.ResponseWriter, r *http.Request){
	if "POST" != r.Method {
		showUploadHtml(w)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil{
		panic(model.CustomException{402,err.Error()})
	}
	id := r.FormValue("id")
	if id == "" {
		panic(model.CustomException{402,"id 不能为空"})
	}
	defer file.Close()
	if !helpers.IsDirExists(config.All().UploadPath+"/"+id) {
		panic(model.CustomException{402,id+"目录不存在！"})
	}
	f,err:=os.Create(config.All().UploadPath+"/"+id+"/"+header.Filename)
	if err != nil{
		panic(model.CustomException{403,id+err.Error()})
	}
	defer f.Close()
	_,err = io.Copy(f,file)
	if err != nil {
		panic(model.CustomException{500,err.Error()})
	}
	//同步返回
	helpers.ResponseJson(200,"上传文件成功!",helpers.EmptyResult,w)
	return
}

func showUploadHtml( w http.ResponseWriter) {
	// 上传页面
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	html := `
<form enctype="multipart/form-data" action="/file/upload" method="POST">
    上传文件: <input name="file" type="file" />
 <input name="id" type="text" />
    <input type="submit" value="Send File" />
</form>
`
	io.WriteString(w, html)
	//helpers.ResponseJson(405,"Method Not Allow",helpers.EmptyResult,w)
}
