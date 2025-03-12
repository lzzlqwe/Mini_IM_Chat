package server

import (
	"net"
	"strings"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
	server  *Server
}

// 监听当前User的channel，一但有消息，就发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.Channel
		this.conn.Write([]byte(msg + "\r \n"))
	}
}

// 用户上线功能
func (this *User) Online() {
	//用户上线，将用户写入到OnlineMap中
	this.server.MapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.MapLock.Unlock()

	//广播当前用户的上线信息
	this.server.BroadCast(this, "already online")
}

// 用户下线功能
func (this *User) Offline() {
	//用户下线，删除OnlineMap中的user信息
	this.server.MapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.MapLock.Unlock()

	this.server.BroadCast(this, "offline")
}

// 给当前用户对应的客户端发送消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg + "\r \n"))
}

// 用户处理消息的业务
func (this *User) DoMeaage(msg string) {
	if msg == "who" { //如果用户发送who指令，则查询当前在线用户
		this.server.MapLock.Lock()
		for _, user := range this.server.OnlineMap { //查询当前在线的用户有哪些
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "on line..."
			this.SendMessage(onlineMsg)
		}
		this.server.MapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" { //修改当前用户名
		//消息格式 rename|jack
		newName := strings.Split(msg, "|")[1]

		//判断name是否以及存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("this name already exist")
		} else {
			this.server.MapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.Name = newName
			this.server.OnlineMap[newName] = this
			this.server.MapLock.Unlock()

			this.SendMessage("you have changed your name to " + newName)
		}

	} else {
		//将得到的消息进行广播
		this.server.BroadCast(this, msg)
	}
}

// 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:    conn.RemoteAddr().String(),
		Addr:    conn.RemoteAddr().String(),
		Channel: make(chan string),
		conn:    conn,
		server:  server,
	}
	//启动goroutine，用于监听当前User的channel消息
	go user.ListenMessage()
	return user
}
