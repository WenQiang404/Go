package main

//
import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在此处创建OnlineMap
	OnlineMap map[string]*User
	//创建读写锁
	mapLock sync.RWMutex
	//消息广播的channel
	Message chan string
}

//创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string, 10),
	}
	return server
}

//start services

func (this *Server) Start() {
	//socket listening
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error", err)
		return
	}
	//close listen socket
	defer Listener.Close()

	//一旦建立连接需要时刻监听channel是否有消息
	//创建一个goroutine来启动Listener方法
	go this.ListenMessage()

	//循环监控是否有连接产生
	for {
		//accept
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Listener error", err)
			continue

		}
		//do handler
		go this.Handler(conn)
	}
}

func (this *Server) Handler(conn net.Conn) {

	//创建一个user
	user := NewUser(conn)
	//当前连接的业务
	// fmt.Println("success to connect!")
	//用户上线
	//因为要向OnlineMap中写入数据，所以需要先上锁
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//广播用户上线的消息
	this.Transfer(user, "the user has logging")

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("conn read err", err)
				return
			}
			if n == 0 {
				this.Transfer(user, "logging out...")
			}

			//提取用户的消息，并去除'\n'
			msg := string(buf[:n-1])

			//广播
			this.Transfer(user, msg)
		}

	}()
	//当前handler阻塞，避免go层死亡
	select {}
}

//创建广播方法
func (this *Server) Transfer(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ":" + msg //广播的信息
	this.Message <- sendMsg                                   //将要广播的消息发送到channel中
}

//创建监听message广播消息的方法，一旦channel中产生消息，就发送给user.channel
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, i := range this.OnlineMap {
			i.Channel <- msg
		}
		this.mapLock.Unlock()
		//将message发送给全部的user

	}
}

//添加onlineUser的map表(无需新创建结构体，绑定到server)
// type OnlineUser struct {
// 	UserName string
// 	UserAddress string
// }

//添加channel管道，广播总消息
