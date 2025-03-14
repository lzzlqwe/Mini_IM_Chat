package client

import (
	"fmt"
	"io"
	"net"
	"os"
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

// 更新用户名
func (this *Client) updateName() bool {
	fmt.Println(">>>>>>>>>>>请输入用户名")
	fmt.Scanln(&this.Name)

	sendMsg := "rename|" + this.Name
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return false
	}
	return true
}

// 处理Server返回的消息，直接显示到终端
func (this *Client) DealResponse() {
	//一但client.conn有数据，就copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, this.conn)
}

// 公聊模式
func (this *Client) PublicChat() {
	fmt.Println(">>>>>>>>>>>请输入聊天内容，exit退出")
	var chatMsg string
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		//消息不为空则发送
		if len(chatMsg) != 0 {
			//发送给服务器
			_, err := this.conn.Write([]byte(chatMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>>>>>>>请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

// 查询在线用户
func (this *Client) SelectUsers() {
	sendMsg := "who"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

// 私聊模式
func (this *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	this.SelectUsers()
	fmt.Println(">>>>>>>>>>>请输入聊天对象[用户名]，exit退出:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>>>>>>>请输入聊天内容，exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//消息不为空则发送
			if len(chatMsg) != 0 {
				//发送给服务器
				sendMsg := "to|" + remoteName + "|" + chatMsg
				_, err := this.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>>>>>>>请输入聊天内容，exit退出:")
			fmt.Scanln(&chatMsg)
		}

		this.SelectUsers()
		fmt.Println(">>>>>>>>>>>请输入聊天对象[用户名]，exit退出:")
		fmt.Scanln(&remoteName)
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
			this.PublicChat()
			break
		case 2:
			//私聊模式
			this.PrivateChat()
			break
		case 3:
			//更新用户名
			this.updateName()
			break
		}
	}
}

func NewClient(ServerIp string, ServerPort int) *Client {
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
