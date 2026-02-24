package client

import (
	"fmt"
	"io"
	"net"
	"os"
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

func (c *Client) DealResponse() {
	_, err := io.Copy(os.Stdout, c.Conn) //io.Copy函数的特殊封装，当读到EOF时不会返回EOF错误，而是正常结束
	if err != nil {
		fmt.Println("\n连接异常断开:", err)
	} else {
		fmt.Println("\n服务器已断开连接。")
	}
	os.Exit(0)
}

func (c *Client) printMenu() {
	fmt.Println("------- 菜单 -------")
	fmt.Println("请选择操作:")
	fmt.Println("1. 公聊功能")
	fmt.Println("2. 私聊功能")
	fmt.Println("3. 更新名称")
	fmt.Println("0. 退出")
}

func (c *Client) Run() error {
	go c.DealResponse()
	c.printMenu()
	for {
		var choice string
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			err := c.PublicChat()
			if err != nil {
				return err
			}
			// 公聊结束后继续显示菜单
			c.printMenu()
		case "2":
			err := c.PrivateChat()
			if err != nil {
				return err
			}
			// 私聊结束后继续显示菜单
			c.printMenu()
		case "3":
			err := c.updateName()
			if err != nil {
				return err
			}
		case "0":
			fmt.Println("退出中...")
			err := c.Conn.Close()
			return err
		default:
			fmt.Println("无效选择，请输入 0-3。")
			c.printMenu()
		}
	}
}

func (c *Client) updateName() error {
	fmt.Print("请输入新的用户名: ")
	var newName string
	fmt.Scanln(&newName)
	c.Name = newName
	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.Conn.Write([]byte(sendMsg))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PublicChat() error {
	fmt.Println("进入公聊模式，输入 'exit' 退出")
	for {
		var msg string
		fmt.Scanln(&msg)
		if msg == "exit" {
			fmt.Println("退出公聊模式")
			break
		}
		sendMsg := msg + "\n"
		_, err := c.Conn.Write([]byte(sendMsg))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) PrivateChat() error {
	fmt.Println("进入私聊模式，输入 'exit' 退出")
	fmt.Print("请输入私聊对象的用户名: ")
	var targetName string
	fmt.Scanln(&targetName)
	for {
		var msg string
		fmt.Scanln(&msg)
		if msg == "exit" {
			fmt.Println("退出私聊模式")
			break
		}
		sendMsg := "to|" + targetName + "|" + msg + "\n"
		_, err := c.Conn.Write([]byte(sendMsg))
		if err != nil {
			return err
		}
	}
	return nil
}
