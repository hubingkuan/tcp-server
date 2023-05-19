package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"server-demo/iface"
	"server-demo/util"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// message中DataLen(uint32 4字节)+Id(uint32 4字节)
	return 8
}

func (dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	//  先将消息长度读取出来 然后根据信息中的data长度进行读取
	dataBuff := bytes.NewReader(binaryData)
	// 只解压head信息 得到datalen和msgID
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}
	// 判断dataLen是否超时允许的最大包长度
	if util.GlobalObject.MaxPackageSize > 0 && msg.dataLen > util.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv!")
	}
	return msg, nil
}
