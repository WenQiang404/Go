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
	flag       int //model
}

func NewClient(serverIp string, serverPort int) *Client {

	//创建客户端对象
	Client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更行用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Print(">>>>>>>please enter the valid number")
		return false
	}

}

func (client *Client) Run() {
	for client.flag != 0 { //循环判断客户端不退出
		for !client.menu() {
		}

		//根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			fmt.Println("public chat")
			break
		case 2:
			fmt.Println("private chat")
			break
		case 3:
			fmt.Println("change your username")
			break
		}

	}
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
	client.Run()

}
