package ziface

import "github.com/gorilla/websocket"

//接口定义
type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行状态
	Serve()
	//添加路由
	SetRouter(router IRouter)
	//返回 连接管理
	GetConnMgr() IConnManager
	//连接之前回调
	SetOnConnStart(func(conn IConnection))
	//关闭之前回调
	SetOnConnStop(func(conn IConnection))
	//调用连接之前
	CallOnConnStart(conn IConnection)
	//调用关闭之前
	CallOnConnStop(conn IConnection)
	//超过最大连接回调
	SetOnMaxConn(func(Conn *websocket.Conn))
	//超过最大连接回调
	CallOMaxConn(conn *websocket.Conn)
}
