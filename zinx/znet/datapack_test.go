package znet

import (
	"encoding/json"
	// "fmt"
	"testing"
	"zinxWebsocket/zinx/znet"
)

//测试结构体
type xiaoxi struct {
	action int
	name   string
}

//只测试datapack封包，拆包
func TestDataPack(t *testing.T) {
	// fmt.Print("TestDataPack")
	var id uint32
	id = 200
	mt := 1
	data := &xiaoxi{}
	data.action = 300
	data.name = "hello,张三"
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(jsonData)
	message := &znet.Message{Id: id, Data: jsonData, MessageType: mt}
	dp := &znet.DataPack{MessageType: mt}

	dataPack, err := dp.Pack(message)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(dataPack)
	outMsg, err := dp.Unpack(mt, dataPack)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(outMsg)
}
