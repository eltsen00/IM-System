package main

import (
	"fmt"

	"github.com/eltsen00/IM-System/server"
)

func main() {
	s := server.NewServer("127.0.0.1", 8888)
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
