package util

import (
	"encoding/json"
	"os"
	"zinx-demo/iface"
)

type GlobalObj struct {
	// 全局server对象
	TcpServer iface.IServer
	// server允许的最大链接
	MaxConn int `json:"max_conn"`
	// tcp端口号
	TcpPort int `json:"tcp_port"`
	// 数据包的最大值
	MaxPackageSize uint32 `json:"max_package_size"`
	// worker池的goroutine数量
	WorkerPoolSize uint32 `json:"worker_pool_size"`
	// 配置单个worker对应的消息队列中任务数量的最大值
	MaxWorkerTaskLen uint32 `json:"max_worker_pool_size"`
	// 监听的ip
	Host string `json:"host"`
	// 服务器名称
	Name string `json:"name"`
	// 服务器版本号
	Version string `json:"version"`
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		MaxConn:        100,
		TcpPort:        8999,
		MaxPackageSize: 4096,
		Host:           "0.0.0.0",
		Name:           "maguahu_server",
		Version:        "v0.1",
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}