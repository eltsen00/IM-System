package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int
	// 用户在线映射表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex // 读写协同锁，保证在线用户映射表的并发安全
	// 消息广播的channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) ListenMessager() {
	for msg := range this.Message {
		this.mapLock.RLock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.RUnlock()
	}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// 处理用户连接
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播用户上线消息
	this.BroadCast(user, "已上线")

	// 接收用户发送的消息
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			// 用户下线
			this.BroadCast(user, "已下线")
			this.mapLock.Lock()
			delete(this.OnlineMap, user.Name)
			this.mapLock.Unlock()
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("Conn Read err:", err)
			return
		}
		msg := string(buf[:n-1])
		this.BroadCast(user, msg)
	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close() // 关闭监听器

	// 启动监听Message的goroutine
	go this.ListenMessager()

	fmt.Printf("服务器已启动，监听 %s:%d\n", this.IP, this.Port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go this.Handler(conn)
	}
}
