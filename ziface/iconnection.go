package ziface

import (
	"net"

	"github.com/gorilla/websocket"
)

//连接管理
type IConnection interface {
	//启动连接，让当前连接，开始工作 此接口不能暴露给前端
	// Start()

	//停止连接，结束当前连接工作
	Stop()

	//获取当前连接的websocket conn
	GetWsConnection() *websocket.Conn

	//获取当前连接的id
	GetConnID() uint32

	//获取连接客户端的状态 ip 端口
	RemoteAddr() net.Addr

	//发送数据，将数据发送给远程客户端（无缓冲）
	SendMsg(data string) error

	//发送数据，将数据发送给远程客户端（有缓冲）
	SendBuffMsg(data string) error

	//发送数据，将数据发送给远程客户端（无缓冲）
	SendByteMsg(data []byte) error

	//发送数据，将数据发送给远程客户端（有缓冲）
	SendBuffByteMsg(data []byte) error

	//设置连接属性
	SetProperty(key string, value interface{})

	//获取连接属性
	GetProperty(key string) (interface{}, error)

	//移除设置属性
	RemoveProperty(key string)

	//得到当前消息类型
	GetMessageType() int

	//是否关闭
	IsClosed() bool
}
