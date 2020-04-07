package main

import (
	"log"
	"zinxWebsocket/zinx/ziface"
	"zinxWebsocket/zinx/znet"
)

//测试路由
type PingRouter struct {
	znet.BaseRouter
}

//处理业务之前
func (br *PingRouter) PreHandle(request ziface.IRequest) {
	// log.Println("test PreHandle")
	err := request.GetConnection().GetWsConnection().WriteMessage(request.GetMessageType(), []byte("ping before"))
	if err != nil {
		log.Println("test Handle err:", err)
	}
}

//处理业务
func (br *PingRouter) Handle(request ziface.IRequest) {
	// log.Println("test Handle")
	err := request.GetConnection().GetWsConnection().WriteMessage(request.GetMessageType(), []byte("ping ping"))
	if err != nil {
		log.Println("test Handle err:", err)
	}
}

//处理业务之后
func (br *PingRouter) PostHandle(request ziface.IRequest) {
	// log.Println("test PostHandle")
	err := request.GetConnection().GetWsConnection().WriteMessage(request.GetMessageType(), []byte("ping after"))
	if err != nil {
		log.Println("test Handle err:", err)
	}
}

func main() {
	//创建一个实例
	s := znet.NewServer()
	//添加路由
	s.AddRouter(&PingRouter{})
	//启动
	s.Serve()
}
