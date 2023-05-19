package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"server-demo/iface"
	"server-demo/util"
)

type Connection struct {
	// server
	TcpServer iface.IServer
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接的ID
	ConnID uint32
	// 当前的链接状态
	isClosed bool
	// 告知当前链接已经退出/停止 channel
	ExitChan chan bool
	// 无缓冲管道 用于读写goroutine之间的通信
	msgChan chan []byte
	// 该链接处理的方法
	MsgHandler iface.IMsgHandler
}

func (c *Connection) StartRead() {
	fmt.Println("[Read Goroutine is running]")
	defer fmt.Println("connID=", c.ConnID, "Read is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		// //  读取客户端的数据到缓冲区
		// buf := make([]byte, util.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err, ", err)
		// 	continue
		// }

		// 创建拆包对象
		dp := NewDataPack()
		// 读取消息头
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTcpConnection(), headData)
		if err != nil {
			fmt.Println("read eof err", err)
			break
		}
		msg, err := dp.Unpack(headData)
		if msg.GetDataLen() > 0 {
			bodyData := make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), bodyData); err != nil {
				fmt.Println("read eof err", err)
				break
			}
			msg.SetData(bodyData)
		}
		// 得到当前conn数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if util.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制
			// 从消息中读取msgID找到绑定的msgID对应的handler处理消息
			go c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Gortine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit]")
	for {
		// 不断的阻塞等待channel的消息 回写给客户端
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error", err)
				return
			}
		case <-c.ExitChan:
			// 代表此时Reader已经退出 writer也要退出
			return

		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID:", c.ConnID)
	// 启动 从当前链接的读取数据业务
	go c.StartRead()
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID:", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	// 关闭socket链接
	c.Conn.Close()
	// 向Writer协程发送信息
	c.ExitChan <- true
	// 连接管理器删除连接
	c.TcpServer.GetConnManager().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	// 将data进行封包
	dp := NewDataPack()
	// 封装数据并且返回二进制的数据
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("pack error msg id=", msgID)
		return errors.New("pack error msg")
	}
	// 将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

// 初始化链接模块的方法
func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
	}
	server.GetConnManager().Add(c)
	return c
}
