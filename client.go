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
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1,公聊模式")
	fmt.Println("2,私聊模式")
	fmt.Println("3,更新用户名")
	fmt.Println("0,退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的数字")
		return false
	}

}

func (client *Client) UpdataName() bool {
	fmt.Println(">>>>>请输入用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remotenName, chatMsg string

	client.SelectUsers()
	fmt.Println("请输入聊天对象的用户名，exit退出")
	fmt.Scanln(&remotenName)

	for remotenName != "exit" {
		fmt.Println(">>>>>>>请输入消息内容")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remotenName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>>>>请输入内容，exit退出")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println("请输入聊天对象的用户名，exit退出")
		fmt.Scanln(&remotenName)
	}
}

func (client *Client) dealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) PublicChat() {

	var chatMsg string

	fmt.Println(">>>>>>>>请输入内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}

		}
		chatMsg = ""
		fmt.Println(">>>>>>>>请输入内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}

		switch client.flag {

		case 1:
			client.PublicChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.UpdataName()
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip地址")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口")
}
func main() {

	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>>链接失败")
		return
	}
	go client.dealResponse()

	fmt.Println(">>>>>>>>>>>>>>>链接服务器成功")

	client.Run()
}
