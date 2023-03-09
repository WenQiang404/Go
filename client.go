package main

import (
	"flag"
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

var serverIp string
var serverPort int

//  ./client -ip 127.0.0.1
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set the server ip address(the default value is 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set the server port(the default value is 8888)")

}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>failed to connect")
		return
	}

	fmt.Println(">>>>>>>>success to connect")

	//启动客户端服务
	select {}
}
