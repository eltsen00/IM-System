package server

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string, 100),
		server: server,
		conn:   conn,
	}
	return user
}

func (this *User) ListenMessage() {
	for msg := range this.C {
		this.conn.Write([]byte(msg + "\n"))
	}
	this.conn.Close()
}

func (this *User) Online() {
	// 处理用户连接
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	// 启动监听当前 User channel 消息的 goroutine
	go this.ListenMessage()
	// 广播用户上线消息
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	// 处理用户断开连接
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	close(this.C)
	// 广播用户下线消息
	this.server.BroadCast(this, "已下线")
}

func (this *User) SendMsg(msg string) {
	if msg == "who" {
		this.server.mapLock.RLock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": 在线..."
			this.C <- onlineMsg
		}
		this.C <- "当前在线用户数量: " + fmt.Sprintf("%d", len(this.server.OnlineMap)) + "\n"
		this.server.mapLock.RUnlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		if _, ok := this.server.OnlineMap[newName]; ok {
			this.C <- "用户名已存在，请重新选择一个用户名\n"
			return
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.C <- "用户名已更改为: " + this.Name + "\n"
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		parts := strings.SplitN(msg, "|", 3)
		if len(parts) != 3 {
			this.C <- "消息格式错误，正确格式: to|用户名|消息内容\n"
			return
		}
		// 1. 获取目标用户名和消息内容
		targetName := parts[1]
		message := parts[2]
		if targetName == this.Name {
			this.C <- "不能给自己发送消息\n"
			return
		}
		if targetName == "" {
			this.C <- "目标用户名不能为空\n"
			return
		}
		if message == "" {
			this.C <- "消息内容不能为空\n"
			return
		}
		// 2. 查找目标用户User对象
		this.server.mapLock.RLock()
		targetUser, ok := this.server.OnlineMap[targetName]
		this.server.mapLock.RUnlock()
		if !ok {
			this.C <- "用户 " + targetName + " 不在线\n"
			return
		}
		// 3. 将消息发送给目标用户
		targetUser.C <- "[" + this.Addr + "]" + this.Name + ": " + message
	} else {
		this.server.BroadCast(this, msg)
	}
}
