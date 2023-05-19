package iface

type IMsgHandler interface {
	// 调度执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)

	// 添加路由
	AddRouter(msgID uint32,router IRouter)

	// 启动worker工作池
	StartWorkerPool()

	SendMsgToTaskQueue(requst IRequest)
}