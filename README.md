# AnCloud-Server

AnCloud-Server 是一个基于 Go 语言开发的 AnCloud 的服务端项目，使用 Gin 框架构建 WebRTC 信令服务。

## 快速开始

### 环境要求

- Go 1.23.2 或更高版本
- Docker（可选，用于容器化部署）

### 安装

#### 克隆项目

```bash
git clone https://github.com/simpleyyt/AnCloud-Server.git
cd AnCloud-Server
```

#### 安装依赖

```bash
go mod download
```

### 编译

#### 本地编译

```bash
# Linux/macOS
go build -o ancloud-server

# Windows
go build -o ancloud-server.exe
```

#### 交叉编译

```bash
# 编译 Linux 版本
GOOS=linux GOARCH=amd64 go build -o ancloud-server-linux-amd64

# 编译 Windows 版本
GOOS=windows GOARCH=amd64 go build -o ancloud-server-windows-amd64.exe

# 编译 macOS 版本
GOOS=darwin GOARCH=amd64 go build -o ancloud-server-darwin-amd64
```

### 运行

#### 直接运行

```bash
# 运行源码
go run main.go

# 运行编译后的程序
./ancloud-server  # Linux/macOS
.\ancloud-server.exe  # Windows
```

#### 使用 Docker 运行

```bash
docker build -t ancloud-server .
docker run -p 8080:8080 ancloud-server
```

## 配置

配置文件位于 `config` 目录下，支持通过 YAML 文件进行配置。

## 许可证

[MIT License](LICENSE)
