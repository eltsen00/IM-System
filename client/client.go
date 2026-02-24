package client

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	return client
}

func (c *Client) Connect() error {
	// 连接服务器
	conn, err := net.Dial("tcp", net.JoinHostPort(c.ServerIp, fmt.Sprintf("%d", c.ServerPort)))
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}
