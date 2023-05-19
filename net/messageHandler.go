package znet

import (
	"fmt"
	"server-demo/iface"
	"server-demo/util"
	"strconv"
)

type MsgHandler struct {
	// 不同的msgID不同的处理方法
	Apis map[uint32]iface.IRouter
	// worker池的worker数量
	WorkerPoolSize uint32
	// 负责Worker取任务的队列
	TaskQueue []chan iface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]iface.IRouter),
		WorkerPoolSize: util.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan iface.IRequest, util.GlobalObject.WorkerPoolSize),
	}
}

// 调度执行对应的Router消息处理方法
func (mg MsgHandler) DoMsgHandler(request iface.IRequest) {
	if handler, ok := mg.Apis[request.GetMsgId()]; !ok {
		fmt.Println("api msgID=", request.GetMsgId(), " is not found")
	} else {
		handler.PreHandle(request)
		handler.Handle(request)
		handler.PostHandle(request)
	}
}

// 为对应的消息ID添加消息处理逻辑
func (mg MsgHandler) AddRouter(msgID uint32, router iface.IRouter) {
	if _, ok := mg.Apis[msgID]; ok {
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID)))
	}
	mg.Apis[msgID] = router
	fmt.Println("add api msgID=", msgID, " success")
}

// 启动一个Worker工作池
func (mh MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan iface.IRequest, util.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workID int, taskQueue chan iface.IRequest) {
	fmt.Println("Worker ID=", workID, " is started")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request iface.IRequest) {
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workerID] <- request
}
