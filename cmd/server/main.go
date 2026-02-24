package main

import (
	"fmt"

	"github.com/eltsen00/IM-System/server"
)

func main() {
	s := server.NewServer("127.0.0.1", 8888)
	if s == nil {
		fmt.Println("Failed to create server.")
		return
	}
	if err := s.Start(); err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	select {}
}
