/**
 *	Author Fanxu(746439274@qq.com)
 */
package main

import (
	"asyncTask/config"
	"os"
	"asyncTask/queue"
	"net"
	"log"
	"asyncTask/websocket"
	"fmt"
)

func main(){
	config.InitConf(os.Args)
	queue.Init()
	go queue.Queue.Consumer()
	host := fmt.Sprintf("%s:%d",config.All().WebSocket.Ip,config.All().WebSocket.Port)
	ln, err := net.Listen("tcp",host)
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept err:", err)
		}
		for {
			go websocket.HandleConnection(conn)
		}
	}
}
