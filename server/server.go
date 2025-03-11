package server

import (
	. "Mini_IM_Chat/user"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户链表
	OnlineMap map[string]*User
	maplock   sync.RWMutex //map锁

	//消息广播的channel
	Message chan string
}

// 用于创建一个Server服务器
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: map[string]*User{},
		Message:   make(chan string), //广播消息的channel
	}
	return server
}

// 监听广播channel的goroutine，一但有消息就转发给全部在线User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		this.maplock.Lock()
		for _, cli := range this.OnlineMap {
			cli.Channel <- msg
		}
		this.maplock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := " [" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//当前连接的业务
	//fmt.Println("连接建立成功！")

	user := NewUser(conn)

	//用户上线，将用户写入到OnlineMap中
	this.maplock.Lock()
	this.OnlineMap[user.Name] = user
	this.maplock.Unlock()

	//广播当前用户的上线信息
	this.BroadCast(user, "already online")

	//当前handle阻塞
	select {}
}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听广播Channel的goroutine
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}

}
