package main

import "net"

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
		C:      make(chan string),
		server: server,
		conn:   conn,
	}
	go user.ListenMessage()
	return user
}

func (this *User) ListenMessage() {
	for msg := range this.C {
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	// 处理用户连接
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播用户上线消息
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	// 处理用户断开连接
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播用户下线消息
	this.server.BroadCast(this, "已下线")
}

func (this *User) SendMsg(msg string) {
	if msg == "who" {
		this.server.mapLock.RLock()
		defer this.server.mapLock.RUnlock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": 在线...\n"
			this.conn.Write([]byte(onlineMsg))
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}
