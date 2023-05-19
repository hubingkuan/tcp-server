package iface

type IRequest interface {
	// 获取当前链接
	GetConnection() IConnection

	// 得到请求的消息数据
	GetData() []byte

	GetMsgId() uint32
}