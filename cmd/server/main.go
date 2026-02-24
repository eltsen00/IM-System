package main

import (
	"flag"
	"fmt"

	"github.com/eltsen00/IM-System/server"
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
	s := server.NewServer(serverIp, serverPort)
	if s == nil {
		fmt.Println("创建服务器失败！")
		return
	}
	if err := s.Start(); err != nil {
		fmt.Println("启动服务器失败:", err)
		return
	}
	select {}
}
