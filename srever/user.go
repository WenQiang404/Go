package main

import (
	"net"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
	server  *Server
}

//创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User { //添加server参数，表示当前用户属于哪个server
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
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

//用户上线的业务
func (this *User) Online() {
	//用户上线
	//因为要向OnlineMap中写入数据，所以需要先上锁
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播用户上线的消息
	this.server.Transfer(this, "the user has logging in")

}

//用户下线的业务
func (this *User) Offline() {
	//用户下线
	//因为要向OnlineMap中删除数据，所以需要先上锁
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播用户上线的消息
	this.server.Transfer(this, "the user has logging out")
}

//用户处理消息
func (this *User) DoMessage(msg string) {
	this.server.Transfer(this, msg)
}
