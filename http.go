/**
 *	Author Fanxu(746439274@qq.com)
 */
package main

import (
	"net/http"
	"log"
	"fmt"
	"os"
	"asyncTask/control"
	"asyncTask/queue"
	"asyncTask/route"
	"asyncTask/config"
	"strings"
)

func main(){
	config.InitConf(os.Args)
	control.Init()
	queue.Init()
	//helpers.InitDbPool()
	//执行队列
	go queue.Queue.Consumer()
	//go queue.ServerMonitor()
	go queue.SystemMonitor()
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/createToken", createToken)
	http.HandleFunc("/",route.Run)
	listen := fmt.Sprintf("%s:%d", config.All().Http.Ip, config.All().Http.Port)
	log.Println("listen:", listen)
	err := http.ListenAndServe(listen,nil)
	if err != nil {
		log.Fatalln(err)
	}
}
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("listen:", "/favicon.ico")
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/png;")
}
func createToken(w http.ResponseWriter, r *http.Request) {
	ip := strings.Replace(r.RequestURI,"/createToken?","",1)
	fmt.Println(ip)
	//s :=  helpers.HlcEncode(ip)
	s :=  "xxxxxxxx"
	w.Write( []byte( `{"code":200,"msg":"ok","data":{"token":"`+s+`"}}` ) )
}