package main

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
	//3.8我自己写的版本
	// if msg != "findall" {
	// 	this.server.Transfer(this, msg)
	// } else {
	// 	for k, _ := range this.server.OnlineMap {
	// 		this.server.Transfer(this, k)
	// 	}

	// }
	//3.8视频里老师写的版本
	if msg == "who" {
		//查询当前在线用户有哪些

		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "online......\n"
			this.SendMessage(onlineMsg)
		}

		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		//判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("this name has been used by another user")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessage("success")

		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMessage("err! plesse use 'to|xxx' \n")
			return
		}
		//根据用户名，得到对方user对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMessage("user not found")
			return
		}
		content := strings.Split(msg, "|")[2]

		if content == "" {
			this.SendMessage("there is no message to send")
			return
		}
		remoteUser.SendMessage(this.Name + " say to you :" + content)

	} else {
		this.server.Transfer(this, msg)
	}

}

//给当前用户发消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}
