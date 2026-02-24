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

func (c *Client) Run() error {
	go c.DealResponse()
	for {
		fmt.Println("Please choose an option:")
		fmt.Println("1. Send Message")
		fmt.Println("2. Receive Message")
		fmt.Println("3. Update Name")
		fmt.Println("0. Exit")
		var choice string
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			//return c.sendMessage()
		case "2":
			//return c.receiveMessage()
		case "3":
			err := c.updateName()
			if err != nil {
				return err
			}
		case "0":
			fmt.Println("Exiting...")
			err := c.Conn.Close()
			return err
		default:
			fmt.Println("Invalid choice, please enter 0-3.")
			fmt.Println()
		}
	}
}

func (c *Client) updateName() error {
	fmt.Print("Enter your new name: ")
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
