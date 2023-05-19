package iface

type IServer interface{
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()

	AddRouter(msgId uint32, router IRouter)

	// 获取连接管理器
	GetConnManager() IConnManager
}