package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func newClient(ServerIp string, ServerPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
	}

	//连接服务端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn
	//返回客户端对象
	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Set the server IP address (default is 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "Set the server port (default is 8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	client := newClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>>连接服务器失败")
		return
	}
	fmt.Println(">>>>>>>>>>>连接服务器成功")

	//启动客户端业务
	for {

	}
}
