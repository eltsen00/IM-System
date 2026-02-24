package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(conn, this)
	user.Online()

	// 用户是否活跃的channel
	isLive := make(chan bool, 1)

	// 使协程同步退出
	ctx, cancel := context.WithCancel(context.Background())

	// 接收用户发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				select {
				case isLive <- true:
				default:
				}
				msg := string(buf[:n-1]) // 去掉末尾的换行符
				user.SendMsg(msg)
			}

			if err != nil {
				if err == io.EOF {
					fmt.Println("用户由客户端主动下线")
				} else {
					fmt.Println("Conn Read err:", err)
				}
				user.Offline()
				cancel() // 取消上下文，通知其他协程退出
				return
			}
		}
	}()
	for {
		select {
		case <-isLive:
			// 当前用户活跃，应该重置定时器
		case <-ctx.Done():
			// 上下文被取消，退出当前协程
			return
		case <-time.After(time.Second * 300):
			// 判断用户是否已经下线
			this.mapLock.RLock()
			_, ok := this.OnlineMap[user.Name]
			this.mapLock.RUnlock()
			if !ok {
				// 用户已经下线，关闭定时器
				return
			}
			// 已经超时，关闭用户连接
			user.C <- "你被踢了,因为你已经300秒没有发消息了"
			user.Offline()
			return
		}
	}
}

func (this *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		return err
	}
	defer listener.Close() // 关闭监听器

	// 启动监听Message的goroutine
	go this.ListenMessager()

	fmt.Printf("服务器已启动，监听 %s:%d\n", this.IP, this.Port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Accept err:", err)
			continue
		}

		go this.Handler(conn)
	}
}
