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
	flag       int //表示当前Client的模式
}

// 获取用户的输入模式
func (this *Client) menu() bool {
	fmt.Println("=======================================================")

	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>>>>>>请输入合法范围内的数字")
		return false
	}
}

// 客户端主业务
func (this *Client) Run() {
	for this.flag != 0 {
		for this.menu() != true {
		}
		switch this.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式连接...")
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式连接...")
			break
		case 3:
			//更新用户名
			fmt.Println("更新用户名...")
			break
		}
	}
}

func newClient(ServerIp string, ServerPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		flag:       -1,
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
	client.Run()
}
