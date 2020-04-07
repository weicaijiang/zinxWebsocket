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
	// log.Println("test Handle")
	// err := request.GetConnection().GetWsConnection().WriteMessage(request.GetMessageType(), []byte("ping ping"))
	// if err != nil {
	// 	log.Println("test Handle err:", err)
	// 	return
	// }

	//base64
	// base64Data, err := base64.StdEncoding.DecodeString(string(request.GetData()))
	// if err != nil {
	// 	log.Println("test Handle DecodeString err:", err)
	// 	return
	// }
	//unmarsh
	msg := &message.Account{}
	err := json.Unmarshal(request.GetData(), msg)
	if err != nil {
		log.Println("test Handle Unmarshal err:", err, " msg:", msg)
		return
	}
	log.Println("test Handle recv server msg:", msg)

	//回写
	responseMsg := &message.Room{Name: "niuniu好好玩", Level: 2, Port: 3}
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

func main() {
	//创建一个实例
	s := znet.NewServer()
	//添加路由
	s.AddRouter(&PingRouter{})
	//启动
	s.Serve()
}
