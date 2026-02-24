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
		fmt.Println("Failed to create client.")
		return
	}
	if err := c.Connect(); err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	fmt.Println("Connected to server successfully.")
	err := c.Run()
	if err != nil {
		fmt.Println("Error running client:", err)
	}
}
