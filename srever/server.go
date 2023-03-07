package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
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
	//当前连接的业务
	fmt.Println("success to connect!")
}
