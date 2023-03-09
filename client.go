package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

//需要一个goroutine持续监听服务器端返回的消息
func (client *Client) DealResponse() {
	//一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

//菜单显示功能
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
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

//1.公聊模式
func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>>please input the message,'exit' to end")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		//send to the server
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>>please input the message,'exit' to end")
		fmt.Scanln(&chatMsg)
	}
}

//2.私聊模式
//2.1查询在线用户
func (client *Client) FindUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err :", err)
		return
	}

}
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.FindUser()
	fmt.Println("please input the user you want to chat Privacy,'exit' to quit")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>please input the message,'exit' to quit")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			//消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>>please input the message,'exit' to end")
			fmt.Scanln(&chatMsg)
		}
		break

	}
}

//3.更新用户名功能
func (client *Client) UpdateName() bool {
	fmt.Println("please input you name")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "/n"

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("connect.Write err", err)
		return false

	}
	return true
}

//启动客户端服务
func (client *Client) Run() {
	for client.flag != 0 { //循环判断客户端不退出
		for !client.menu() {
		}

		//根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
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

	//单独开启一个goroutine

	fmt.Println(">>>>>>>>success to connect")
	go client.DealResponse()
	//启动客户端服务
	client.Run()

}
