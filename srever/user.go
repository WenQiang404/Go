package main

import (
	"net"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
}

//创建一个用户的API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
	}

	//启动监听channel的goroutine
	go user.ListenMessage()
	return user
}

//监听当前每个user中的channel方法,若有消息直接发送给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.Channel

		this.conn.Write([]byte(msg + "\n"))
	}
}
