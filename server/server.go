package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户链表
	OnlineMap map[string]*User
	MapLock   sync.RWMutex //map锁

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

		this.MapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.Channel <- msg
		}
		this.MapLock.Unlock()
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

	user := NewUser(conn, this)

	//用户上线
	user.Online()

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 { //客户端断开连接，用户下线
				user.Offline()
				return
			}

			if err != nil {
				fmt.Println("conn Read err:", err)
				return
			}

			//提取用户发送的消息(去除'\n')
			//msg := string(buf[:n-1])
			msg := string(buf[:n])
			//针对message进行处理
			user.DoMeaage(msg)
		}
	}()

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
