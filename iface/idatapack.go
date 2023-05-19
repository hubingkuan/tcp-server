package iface


/**
封包 拆包模块
面向tcp链接的数据流 自定义消息格式解决粘包问题
*/
type IDataPack interface {
	// 获取包的消息长度
	GetHeadLen() uint32
	// 封包
	Pack(msg IMessage)([]byte,error)
	// 拆包
	Unpack([]byte) (IMessage,error)
}