package znet

import (
	// "encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"zinxWebsocket/zinx/utils"
	"zinxWebsocket/zinx/ziface"
)

//封包，拆包模块
type DataPack struct{}

//拆包，封包
func NewDataPack() *DataPack {
	return &DataPack{}
}

//封包
func (dp *DataPack) Pack(imsg ziface.IMessage) (string, error) {

	//第一步组成这样的数据 {Id:1, Data:"hello",MessageType:1}
	msg := &Message{Id: imsg.GetMsgId(), MessageType: imsg.GetMessageType(), Data: imsg.GetData()}
	// log.Println("DataPack Pack msg:", msg)
	// log.Println("DataPack Pack msg.data1:", msg.Data)
	// log.Println("DataPack Pack msg.data2:", string(msg.Data))
	jsonData, err := json.Marshal(msg)
	// log.Println("DataPack Pack jsonData1:", jsonData)
	// log.Println("DataPack Pack jsonDat2:", string(jsonData))
	if err != nil {
		log.Println("DataPack Pack Marshal err:", err, " msg:", msg)
		return "", errors.New("DataPack Pack Marshal err")
	}
	//第二步base64加密
	// base64Data := base64.StdEncoding.EncodeToString(jsonData)

	//判断下长度是否超出大小
	if len(jsonData) > int(utils.GlobalObject.MaxPackageSize) {
		log.Println("DataPack Pack len err jsonData:", jsonData, " len:", len(jsonData))
		return "", errors.New("DataPack Pack msg len big then MaxPackageSize err")
	}
	// log.Println("DataPack Pack jsonData:", string(jsonData))
	return string(jsonData), nil
}

//拆包
func (dp *DataPack) Unpack(messageType int, data string) (ziface.IMessage, error) {

	//暂时直接返回，后续如果有的话，多个消息之间使用 # 号隔开
	// msg := &Message{MessageType: messageType, Data: data}

	//base64解密
	// jsonData, err := base64.StdEncoding.DecodeString(string(data))
	// if err != nil {
	// 	log.Println("DataPack Unpack DecodeString err:", err, " data:", string(data))
	// 	return nil, errors.New("DataPack Unpack DecodeString data err")
	// }
	// log.Println("DataPack Unpack DecodeString data1:",string(data))

	//json解析
	imsg := &Message{}
	err := json.Unmarshal([]byte(data), imsg)
	if err != nil {
		log.Println("DataPack Unpack Unmarshal err:", err, " data:", string(data), " imsg:", imsg)
		return nil, errors.New("DataPack Unpack Unmarshal data err")
	}
	// msg := &Message{Id: d.Id, MessageType: messageType, Data: d.Data}
	return imsg, nil
}
