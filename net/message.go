package znet

type Message struct {
	id      uint32
	dataLen uint32
	data    []byte
}

func NewMsgPackage(id uint32,data []byte) *Message{
	return &Message{
		id: id,
		dataLen: uint32(len(data)),
		data: data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.id
}

func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetMsgId(id uint32) {
	m.id = id
}

func (m *Message) SetData(data []byte) {
	m.data = data
}

func (m *Message) SetDataLen(len uint32) {
	m.dataLen = len
}