package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {

	//创建客户端对象
	Client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("Error connecting...", err)
		return nil
	}

	Client.conn = conn

	return Client

}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>>failed to connect")
		return
	}

	fmt.Println(">>>>>>>>success to connect")

	//启动客户端服务
	select {}
}
