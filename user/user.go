package user

import "net"

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
}

// 监听当前User的channel，一但有消息，就发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.Channel
		this.conn.Write([]byte(msg + "\r \n"))
	}
}

// 创建一个用户
func NewUser(conn net.Conn) *User {
	user := &User{
		Name:    conn.RemoteAddr().String(),
		Addr:    conn.RemoteAddr().String(),
		Channel: make(chan string),
		conn:    conn,
	}
	//启动goroutine，用于监听当前User的channel消息
	go user.ListenMessage()
	return user
}
