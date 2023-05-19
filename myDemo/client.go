package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx-demo/net"
)

func main() {
	fmt.Println("client start...")
	// 直接链接远程服务器 得到conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("client connect err", err)
		return
	}
	// 写数据
	for {
		// 发送封包的message消息 msgID
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("你好服务器")))
		if err != nil {
			fmt.Println("pack error:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error:", err)
			return
		}
		// 解析服务器的回写消息
		headBuf := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headBuf); err != nil {
			fmt.Println("read head err", err)
			break
		}
		message, err := dp.Unpack(headBuf)
		if err != nil {
			fmt.Println("unpack err", err)
			break
		}
		if message.GetDataLen()>0{
			dataBuf := make([]byte, message.GetDataLen())
			if _, err := io.ReadFull(conn, dataBuf);err!=nil{
				fmt.Println("read msg data err",err)
				break
			}
			fmt.Println("recv server msg:ID=",message.GetMsgId()," len:",message.GetDataLen(),"data:",string(message.GetData()))
		}
	}
}