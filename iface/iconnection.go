package iface

import (
	"net"
)

type IConnection interface {
	// 启动链接  让当前的链接准备开始工作
	Start()
	// 停止链接 结束当前链接工作
	Stop()
	// 获取当前链接的绑定socket conn
	GetTcpConnection() *net.TCPConn
	// 获取当前链接模块的链接 ID
	GetConnID() uint32
	// 获取远程客户端的tcp状态 IP+端口
	RemoteAddr() net.Addr
	// 发送数据给 远程的客户端
	Send(msgID uint32,data []byte) error
}

// 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error