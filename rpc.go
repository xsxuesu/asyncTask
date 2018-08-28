/**
 *	Author Fanxu(746439274@qq.com)
 */
package main

import (
	"log"
	"fmt"
	"asyncTask/control"
	"asyncTask/config"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"asyncTask/queue"
)

func main(){
	config.InitConf(os.Args)
	queue.Init()
	go queue.Queue.Consumer()
	listenHost := fmt.Sprintf("%s:%d", config.All().Rpc.Ip, config.All().Rpc.Port)
	listen, err := net.Listen("tcp", listenHost)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer listen.Close()
	host := rpc.NewServer()
	err = host.Register(control.RpcServer{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("listening:", listenHost)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err.Error())
			continue
		}
		go host.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
