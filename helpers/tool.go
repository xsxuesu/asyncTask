/**
 *	Author Fanxu(746439274@qq.com)
 */

package helpers

import (
	"strings"
	"net/http"
	"encoding/json"
	"fmt"
	"net/rpc/jsonrpc"
	"log"
	"net/rpc"
	"asyncTask/config"
	"time"
	"os"
	"os/exec"
	"bytes"
)
type JsonResult struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data map[string]interface{} `json:"data"`
}
var EmptyResult = map[string] interface{}{}
/**
 * 字符串首字母转化为大写 ios_bbbbbbbb -> IosBbbbbbbbb
 */
func StrFirstToUpper(str string) string {
	temp := strings.Split(str, "_")
	var upperStr string
	for y := 0; y < len(temp); y++ {
		vv := []rune(temp[y])
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				vv[i] -= 32
				upperStr += string(vv[i]) // + string(vv[i+1])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	return upperStr
}

func ResponseError(w http.ResponseWriter){
	w.WriteHeader(404)
	w.Write([]byte(`{"code":404,"msg":"404 Method Not Found","data":{}}`))
	return
}

func ResponseJson( code int , msg string , data map[string] interface{} , w http.ResponseWriter ) {
	//允许跨域设置
	w.Header().Set("Content-Type","application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Headers","Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token,Authorization")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Request-Method","GET,HEAD,PUT,PATCH,POST,DELETE")
	var returnData JsonResult
	returnData.Code = code
	returnData.Msg = msg
	returnData.Data = data
	jsonData,err := json.Marshal(returnData)
	if err != nil{
		fmt. Println ( "error:" , err )
		w.Write([]byte(`{"code":500,"msg":"请检查返回数据格式是否为标准interface","data":{}}`))
		return
	}
	log.Printf(string(jsonData))
	w.Write(jsonData)
	return
}

func LoginRequired( w http.ResponseWriter, r *http.Request ) string{
	var resp bool
	globalHost := fmt.Sprintf("%s:%d",config.All().Session.Ip,config.All().Session.Port)
	global := RPCServer(globalHost)
	if global == nil {
		return "error"
	}
	defer global.Close()
	err := global.Call("Remote.SessionValidation", r.Header, &resp)
	if err != nil {
		log.Println("rpc error:", err.Error())
		return "error"
	}
	if !resp {
		return "error"
	}
	return "ok"
}

func RPCServer(domain_port string) *rpc.Client {
	global, err := jsonrpc.Dial("tcp", domain_port)
	if err != nil {
		log.Println("tcp:", domain_port,err.Error())
		return nil
	}
	return global
}

func RPCConn() *rpc.Client {
	serverCount := len(config.All().RpcServers)
	index := time.Now().Unix() % int64(serverCount)
	for {
		server := fmt.Sprintf("%s:%d", config.All().RpcServers[index].Ip, config.All().RpcServers[index].Port)
		log.Println(server)
		client, err := jsonrpc.Dial("tcp", server)
		if err != nil {
			log.Println(err.Error())
			index++
			if int(index) == serverCount {
				index ^= index
			}
			continue
		}
		return client
	}
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	return false
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func ExecShell(s string) (string, error){
	log.Println("执行命令：",s)
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	var err1 bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err1
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	if err != nil{
		log.Println(err)
		return out.String()+err1.String(),err
	}
	return out.String()+err1.String(), err
}

func CheckAuth(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	ips := strings.Split(r.RemoteAddr,":")
	if len(ips) < 2 {
		fmt.Println("来源IP：",ips[0])
		return false
	}
	if auth == "" {
		fmt.Println("来源IP：",ips[0])
		return false
	}
	str := HlcDecode(auth)
	//匹配远程地址
	if str != ips[0] {
		fmt.Println("来源IP：",ips[0])
		fmt.Println("解密IP：",str)
		return false
	}
	//var right = false
	//for i:=0 ; i < len(config.All().WhiteIps) ;i++  {
	//	if config.All().WhiteIps[i] == str {
	//		right =  true
	//		return true
	//	}
	//}
	//if !right {
	//	fmt.Println("IP未在白名单中：",str)
	//}
	return true
}

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func HasElement(arr []string,val string) bool {
	for _,v := range arr  {
		if val == v{
			return true
		}
	}
	return false
}
