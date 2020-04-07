package main

import "zinxWebsocket/zinx/znet"

func main() {
	//创建一个实例
	s := znet.NewServer("zinx websocket v0.1")
	//启动
	s.Serve()
}
