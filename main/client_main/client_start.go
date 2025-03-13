package main

import (
	"Mini_IM_Chat/client"
	"flag"
	"fmt"
)

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Set the server IP address (default is 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "Set the server port (default is 8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	client := client.NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>>连接服务器失败")
		return
	}

	//单独开启一个goroutine去处理server的回执信息
	go client.DealResponse()

	fmt.Println(">>>>>>>>>>>连接服务器成功")

	//启动客户端业务
	client.Run()
}
