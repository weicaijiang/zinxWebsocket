package main

import (
	// "encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"time"
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
	err := json.Unmarshal([]byte(request.GetData()), msg)
	if err != nil {
		log.Println("PingRouter Handle Unmarshal err:", err, " msg:", msg)
		// return
		//不能解析，就直接输出json
		log.Println("PingRouter Handle recv server msg:", string(request.GetData()))
	} else {
		log.Println("PingRouter Handle recv server msg:", msg)
	}
	// log.Println("PingRouter Handle recv server msg:", string(request.GetData()))

	//回写
	responseMsg := &message.Room{Name: "偶是服务器回的包", Level: 2, Port: 3}
	jsonData, err := json.Marshal(responseMsg)
	if err != nil {
		log.Println("PingRouter handler Pack Marshal err:", err, " responseMsg:", responseMsg)
		return
	}
	err = request.GetConnection().SendBuffMsg(string(jsonData))
	if err != nil {
		log.Println("PingRouter Handle WriteMessage err:", err)
		return
	}
	//再发下时间
	request.GetConnection().SendBuffMsg("服务器unix时间:" + strconv.Itoa(int(time.Now().Unix())))
}

//回写数据
type RepeatRouter struct {
	znet.BaseRouter
}

func (br *RepeatRouter) Handle(request ziface.IRequest) {
	//直接取数据回写
	err := request.GetConnection().SendMsg(request.GetData())
	if err != nil {
		log.Println("RepeatRouter Handle WriteMessage err:", err)
		return
	}
	log.Println("RepeatRouter Handle receive from msg:", request.GetData())
}

//回调之后
func DoConectionBegin(conn ziface.IConnection) {
	log.Println("DoConectionBegin is called connid:", conn.GetConnID())
	conn.SendMsg("我在连接开始后的第一个消息")
	conn.SetProperty("haohao", "我爱游戏")
	conn.SetProperty("age", "20")
}

func DoConectionEnd(conn ziface.IConnection) {
	log.Println("DoConectionEnd is called connid:", conn.GetConnID())
	log.Println("我在连接关闭后的最后一条消息")
	// conn.SendMsg(1,3,"conn stop")
	value, err := conn.GetProperty("haohao")
	if err == nil {
		log.Println("DoConectionEnd haohao:", value)
	}
	value, err = conn.GetProperty("age")
	if err == nil {
		log.Println("DoConectionEnd age:", value)
	}
}

func main() {
	//创建一个实例
	s := znet.NewServer()

	//注意连接回调
	s.SetOnConnStart(DoConectionBegin)
	s.SetOnConnStop(DoConectionEnd)

	//回写测试
	s.SetRouter(&RepeatRouter{})

	//返回消息测试
	s.SetRouter(&PingRouter{})

	//启动
	s.Serve()
}
