package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx-demo/iface"
)

type ConnManager struct {
	// 管理连接的集合
	connections map[uint32]iface.IConnection
	// 保护连接集合的读写锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn iface.IConnection) {
	// 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Lock()

	// 将conn加入到connMagager中
	connMgr.connections[conn.GetConnID()]=conn
	fmt.Println("connID:",conn.GetConnID()," add to ConnManager success,conn num=",connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn iface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Lock()
	delete(connMgr.connections,conn.GetConnID())
	fmt.Println("connID:",conn.GetConnID()," remove from ConnManager success,conn num=",connMgr.Len())

}

// 根据ConnID获取链接
func (connMgr *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn,ok := connMgr.connections[connID];ok {
		return conn,nil
	}else{
		return nil,errors.New("connection not found")
	}
}

// 获取当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Lock()
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections,connID)
	}
	fmt.Println("clear all connection succ!,conne num=",connMgr.Len())
}