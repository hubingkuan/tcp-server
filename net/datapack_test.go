package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("listen err", err)
		return
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("Accept err", err)
				return
			}
			go func(conn net.Conn) {
				// 拆包
				datapack := NewDataPack()
				for {
					headData := make([]byte, datapack.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head data err", err)
						return
					}
					// 先读取消息长度
					iMessage, err := datapack.Unpack(headData)
					if err != nil {
						fmt.Println("unpack err", err)
						return
					}

					if iMessage.GetDataLen() > 0 {
						bodyData := make([]byte, iMessage.GetDataLen())
						_, err := io.ReadFull(conn, bodyData)
						if err != nil {
							fmt.Println("read head data err", err)
							return
						}
						fmt.Println("recv msgID:", iMessage.GetMsgId(), "msg:", string(bodyData))
					}
				}
			}(conn)
		}
	}()

	dial, _ := net.Dial("tcp", "127.0.0.1:8999")
	pack := NewDataPack()
	msg1 := &Message{
		id:      1,
		dataLen: 7,
		data:    []byte{'m', 'a', 'g', 'u', 'a', 'h', 'u'},
	}
	msg2 := &Message{
		id:      2,
		dataLen: 7,
		data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	bytes1, err := pack.Pack(msg1)
	if err != nil {
		fmt.Println("pack msg err", err)
		return
	}
	bytes2, err := pack.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg err", err)
		return
	}
	// 模拟粘包情况 将2个消息组合到一起发送
	bytes1 = append(bytes1, bytes2...)
	dial.Write(bytes1)
	select {}

}