package znet

import "zinxWebsocket/zinx/ziface"

//实现router，先嵌入基类
type BaseRouter struct {
}

//处理业务之前
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

//处理业务
func (br *BaseRouter) Handle(request ziface.IRequest) {}

//处理业务之后
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}

//路由
type Router struct {
	BaseRouter
}

//创建实例
// func NewRouter()
