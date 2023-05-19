package znet

import (
	"fmt"
	"net"
	"zinx-demo/util"
	"zinx-demo/iface"
)

// IServer实现类
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
	// 当前server的消息管理  msgID和对应的handler处理
	MsgHandler iface.IMsgHandler
	// 连接管理器
	ConnManager iface.IConnManager
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[zinx]Server Name:%s,ip:%s,port:%d,Maxconn:%d,MaxPackageSize:%d is staring\n", s.Name, s.IP, s.Port, util.GlobalObject.MaxConn, util.GlobalObject.MaxPackageSize)

	// 为了避免AcceptTCP阻塞  异步启动
	go func() {
		// 开启worker工作池
		s.MsgHandler.StartWorkerPool()
		// 获取一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		// 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start  server success, ", s.Name, " Listenning...")
		var cid uint32 = 0
		// 阻塞的等待客户端链接  处理客户端业务(读写)
		for {
			// 如果有客户端链接过来 阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err,", err)
				continue
			}
			// 判断连接是否已满
			if s.ConnManager.Len() > util.GlobalObject.MaxConn {
				fmt.Println("服务端连接已满")
				conn.Close()
				continue
			}
			// 启动链接
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("server name:", s.Name, " stop ")
	s.ConnManager.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	s.Start()

	// 这里由于start是异步启动服务器的 不会阻塞  需要在这里将启动流程阻塞住  不然整个流程就立马结束了
	select {}
}

func (s *Server) AddRouter(msgId uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("add router success")
}

func NewServer() iface.IServer {
	s := &Server{
		Name: util.GlobalObject.Name,
		// 网络有“tcp”、“tcp4”（仅IPv4）、“tcp6”（仅IPv6）、“udp”、“udp4”（仅IPv4）、“udp6”（只IPv6）、“ip”、“ip4”（只IPv4）、《ip6》（只IPv6）、“unix”、“unixgram”和“unixpacket”
		IPVersion:   "tcp4",
		IP:          util.GlobalObject.Host,
		Port:        util.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}

func (s *Server)GetConnManager() iface.IConnManager{
	return s.ConnManager
}