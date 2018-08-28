/**
 *	Author Fanxu(746439274@qq.com)
 */
package main

import (
	"fmt"
	"os"
	"bytes"
	"mime/multipart"
	"io"
	"net/http"
	"log"
	"asyncTask/helpers"
	"strconv"
	"strings"
)

func main(){
	//config.InitConf(os.Args)
	//client := helpers.RPCConn()
	//defer client.Close()
	//var resp map[string]interface{}
	//args := make(map[string]map[string][]string)
	//fmt.Printf("RpcServer.UserIndex")
	//err := client.Call("RpcServer.UserIndex", args, &resp)
	//if err != nil {
	//	fmt.Printf(err.Error())
	//}
	//respJSON, err := json.Marshal(resp)
	//if err != nil {
	//	fmt.Printf(err.Error())
	//}
	//fmt.Printf(string(respJSON))
	cmd := fmt.Sprintf("ls %s -l |grep '^d'|wc -l","/tmp/")
	s,_ := helpers.ExecShell(cmd)
	fmt.Println(s)
	dirCount,err:= strconv.Atoi(strings.Replace(s,"\n","",1))
	if err != nil {
		fmt.Println(err)
	}
	if dirCount >0 {
		fmt.Println("test")
	}
	fmt.Println(dirCount)
}

func Upload() (err error) {
	// Create buffer
	buf := new(bytes.Buffer) // caveat IMO dont use this for large files, \
	// create a tmpfile and assemble your multipart from there (not tested)
	w := multipart.NewWriter(buf)
	// Create file field
	fw, err := w.CreateFormFile("file", "1.sql") //这里的file很重要，必须和服务器端的FormFile一致
	if err != nil {
		fmt.Println("c")
		return err
	}
	fd, err := os.Open("C:/Users/Jiake/Desktop/work/1.sql")
	if err != nil {
		fmt.Println("d")
		return err
	}
	defer fd.Close()
	// Write file field from file to upload
	_, err = io.Copy(fw, fd)
	if err != nil {
		fmt.Println("e")
		log.Println(err)
		return err
	}
	ff, err := w.CreateFormField("id", ) //这里的file很重要，必须和服务器端的FormFile一致
	ff.Write([]byte("asdasdsad"))
	// Important if you do not close the multipart writer you will not have a
	// terminating boundry
	w.Close()
	req, err := http.NewRequest("POST","http://127.0.0.1:8080/file/upload", buf)
	if err != nil {
		fmt.Println("f")
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("g")
		log.Println(err)
		return err
	}
	io.Copy(os.Stderr, res.Body) // Replace this with Status.Code check
	fmt.Println("h")
	return err
}
