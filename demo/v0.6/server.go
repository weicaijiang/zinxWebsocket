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
	dp := znet.NewDataPack()
	msg1 := znet.NewMessage(10, 1, jsonData)
	sendMsg, err := dp.Pack(msg1)

	err = request.GetConnection().GetWsConnection().WriteMessage(request.GetMessageType(), sendMsg)
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
	dp := znet.NewDataPack()
	msg1 := znet.NewMessage(20, 1, jsonData)
	sendMsg, err := dp.Pack(msg1)

	err = request.GetConnection().GetWsConnection().WriteMessage(request.GetMessageType(), sendMsg)
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
	//启动
	s.Serve()
}
