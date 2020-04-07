package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"zinxWebsocket/zinx/ziface"
)

//存储一些数据供全局使用
type GlobalObj struct {
	//当前zinx实例对象
	WsServer ziface.IServer
	//服务器名字
	Name string
	//允许最大连接人数
	MaxConn int
	//当前数据包最大值
	MaxPackageSize uint32

	//类型 ws,wss
	Scheme string
	//连接地址
	Host string
	//端口
	Port uint32
	//子协议
	Path string

	//工作池大小 一般是cpu大小
	WorkerPoolSize uint32
	//一个工作池处理消息的最大数量
	MaxWorkerTaskLen uint32 
}

var GlobalObject *GlobalObj

//重装加载配置
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		log.Println("globalobj reload ReadFile err:", err)
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		log.Println("globalobj reload Unmarshal err:", err)
		return
	}
}

//初始化
func init() {
	GlobalObject = &GlobalObj{
		Name:           "zinx websocket",
		Scheme:         "ws",
		Host:           "0.0.0.0",
		Port:           8999,
		Path:           "",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 4,
		MaxWorkerTaskLen: 1024,
	}

	//从conf加载数据
	GlobalObject.Reload()
}
