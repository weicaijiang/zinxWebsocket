package main

import (
	// "encoding/base64"
	"encoding/json"
	"log"
	"zinxWebsocket/demo/message"
	"zinxWebsocket/zinx/ziface"
	"zinxWebsocket/zinx/znet"
)

//测试路由
type PingRouter struct {
	znet.BaseRouter
}

//处理业务
func (br *PingRouter) Handle(request ziface.IRequest) {
	//unmarsh
	msg := &message.Account{}
	err := json.Unmarshal(request.GetData(), msg)
	if err != nil {
		log.Println("test Handle Unmarshal err:", err, " msg:", msg)
		return
	}
	log.Println("test Handle recv server msg:", msg)

	//回写
	responseMsg := &message.Room{Name: "偶是第一个服务器回的包", Level: 2, Port: 3}
	jsonData, err := json.Marshal(responseMsg)
	if err != nil {
		log.Println("DataPack Pack Marshal err:", err, " responseMsg:", responseMsg)
		return
	}
	err = request.GetConnection().SendMsg(1, 10, jsonData)
	if err != nil {
		log.Println("test Handle WriteMessage err:", err)
		return
	}
}

type SecondRouter struct {
	znet.BaseRouter
}

//处理业务
func (br *SecondRouter) Handle(request ziface.IRequest) {
	//unmarsh
	msg := &message.Account{}
	err := json.Unmarshal(request.GetData(), msg)
	if err != nil {
		log.Println("test Handle Unmarshal err:", err, " msg:", msg)
		return
	}
	log.Println("test Handle recv server msg:", msg)

	//回写
	responseMsg := &message.Room{Name: "偶是第二个服务器回的包哦", Level: 2, Port: 3}
	jsonData, err := json.Marshal(responseMsg)
	if err != nil {
		log.Println("DataPack Pack Marshal err:", err, " responseMsg:", responseMsg)
		return
	}
	err = request.GetConnection().SendMsg(1, 20, jsonData)
	if err != nil {
		log.Println("test Handle WriteMessage err:", err)
		return
	}
}

//回写数据
type RepeatRouter struct {
	znet.BaseRouter
}

func (br *RepeatRouter) Handle(request ziface.IRequest) {
	//直接取数据回写
	err := request.GetConnection().SendMsg(1, 30, request.GetData())
	if err != nil {
		log.Println("test Handle WriteMessage err:", err)
		return
	}
}

func main() {
	//创建一个实例
	s := znet.NewServer()
	//添加路由 收到msgid = 1 返回 10
	s.AddRouter(1, &PingRouter{})
	//收到 msgid=2 返回20
	s.AddRouter(2, &SecondRouter{})
	//收到 msgid=2 返回20
	s.AddRouter(3, &RepeatRouter{})

	//启动
	s.Serve()
}
