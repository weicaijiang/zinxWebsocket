package znet

import "zinxWebsocket/ziface"
import "log"

//实现router，先嵌入基类
type BaseRouter struct {
}

//处理业务之前
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

//处理业务
func (br *BaseRouter) Handle(request ziface.IRequest) {
	//主业务不给路由，就提示一条输出信息
	log.Println("Handle msg:", request.GetData())
}

//处理业务之后
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}

//路由
type Router struct {
	BaseRouter
}
