package main

import (
	"server-demo/iface"
	"server-demo/net"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pingRouter *PingRouter) Handle(request iface.IRequest) {
	request.GetConnection().Send(1, []byte("你好 客户端"))
}

func main() {
	myServer := znet.NewServer()
	myServer.AddRouter(0, &PingRouter{})
	myServer.Serve()
}
