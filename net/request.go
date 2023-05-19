package znet

import (
	"server-demo/iface"
)

type Request struct {
	// 已经和客户端建立好的链接
	conn iface.IConnection

	msg iface.IMessage
}

func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

// 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
