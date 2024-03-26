package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 启动服务器的接口
func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		//fmt.Println("net.Listen err", err)
		return
	}

	defer listener.Close()

	go this.ListenMessager()
	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
		}
		//mt.Println("accept + 1")
		go this.Handler(conn)
		//fmt.Println("accept + 2")
		//nc 127.0.0.1 8888
	}
}

func (this *Server) Handler(conn net.Conn) {

	//fmt.Println("链接建立成功")
	user := NewUser(conn, this)
	//this.mapLock.Lock()
	//this.OnlineMap[user.Name] = user
	//this.mapLock.Unlock()
	//this.BroadCast(user, "已上线")
	user.Online()

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户的消息
			msg := string(buf[:n-1])
			//广播用户消息
			user.DoMessage(msg)
		}
	}()
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	//fmt.Println("accept + 4")
	this.Message <- sendMsg
}

func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}
