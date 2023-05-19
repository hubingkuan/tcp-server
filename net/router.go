package znet

import (
	"server-demo/iface"
)

type BaseRouter struct {
}

func (r *BaseRouter) PreHandle(request iface.IRequest) {

}

func (r *BaseRouter) Handle(request iface.IRequest) {

}

func (r *BaseRouter) PostHandle(request iface.IRequest) {

}
