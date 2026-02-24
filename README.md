# IM-System

一个基于 Go 语言（Golang）开发的简易即时通讯系统。采用 TCP 协议进行通信，支持多人在线聊天、私聊、昵称修改等功能。

## 功能特性

- **TCP 通信**：基于 TCP 协议的稳定连接。
- **消息广播（公聊）**：用户发送的消息会被广播给所有在线用户。
- **私聊功能**：支持点对点私聊消息发送。
- **用户管理**：
    - 用户上线/下线广播通知。
    - 支持在线修改用户名。
    - 在线用户列表维护。
- **超时机制**：服务端自动检测用户活跃状态，若用户超过 300 秒未发送消息，将自动断开连接（踢出）。

## 目录结构

```
IM-System/
├── client/          # 客户端核心逻辑
├── cmd/             # 程序入口
│   ├── client/      # 客户端入口 (main.go)
│   └── server/      # 服务端入口 (main.go)
├── server/          # 服务端核心逻辑
├── go.mod           # Go 模块定义
└── README.md        # 项目说明文档
```

## 环境要求

- Go 1.25.6 或更高版本

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/eltsen00/IM-System.git
cd IM-System
```

### 2. 运行服务端

服务端默认监听 `127.0.0.1:8888`。

```bash
go run cmd/server/main.go
```

你可以通过命令行参数指定 IP 和端口：

```bash
go run cmd/server/main.go -ip 0.0.0.0 -port 8888
```

### 3. 运行客户端

启动客户端连接到服务端。

```bash
go run cmd/client/main.go
```

如果服务端不在默认地址，请指定连接参数：

```bash
go run cmd/client/main.go -ip 127.0.0.1 -port 8888
```

## 使用指南

客户端启动成功后，将显示交互式菜单：

1.  **公聊功能**：进入群聊模式，发送的消息将被所有在线用户收到。
2.  **私聊功能**：可以指定用户名进行一对一聊天。
3.  **更新名称**：修改当前用户的显示名称。
0.  **退出**：断开连接并退出程序。

## ⚙️ 核心技术

- **语言**：Go (Golang)
- **网络**：TCP Socket
- **并发**：Goroutine, Channel, Sync.RWMutex
