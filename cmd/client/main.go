package main

import (
	"flag"
	"fmt"

	"github.com/eltsen00/IM-System/client"
)

var (
	serverIp   string
	serverPort int
)

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Server IP address")
	flag.IntVar(&serverPort, "port", 8888, "Server port")
	flag.Parse() // Parse 的意思是作语法分析
}

func main() {
	c := client.NewClient(serverIp, serverPort)
	if c == nil {
		fmt.Println("连接服务器失败！")
		return
	}
	if err := c.Connect(); err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	fmt.Println("连接服务器成功。")
	err := c.Run()
	if err != nil {
		fmt.Println("运行客户端时出错:", err)
	}
}
