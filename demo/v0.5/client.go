package main

import (
	// "encoding/base64"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"

	// "strconv"

	// "time"
	"zinxWebsocket/demo/message"
	"zinxWebsocket/zinx/znet"

	"github.com/gorilla/websocket"
)

var max = 1

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//	Path: "/echo"
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8999"}
	log.Println("connecting to ", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial err:", err)
		return
	}
	defer conn.Close()

	go timeWriter(conn)

	i := 0
	for {
		log.Println("main ReadMessage start")
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read err:", err)
			return
		}
		log.Println("main ReadMessage read server message", string(msg[:]))
		//解包
		dp := znet.DataPack{}
		recvMsg , err := dp.Unpack(mt,msg)
		if err != nil {
			log.Fatal("main ReadMessage Unpack err:", err)
			return
		}
		log.Println(recvMsg)
		room := &message.Room{}
		err = json.Unmarshal(recvMsg.GetData(),room)
		if err != nil {
			log.Fatal("main ReadMessage Unmarshal err:", err)
			return
		}
		log.Println(room)
		// break
		i++
		if i > max {
			break
		}
	}

}

func timeWriter(conn *websocket.Conn) {
	var i = 0
	for {
		// log.Println("WriteMessage start timeWriter i = ", i)

		//发第一个消息
		msg := &message.Account{Name: "hello,张三", Age: i, Passwd: "123456"}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Println("client timeWriter Marshal err:", err, " msg:", msg)
			break
		}
		log.Println("client timeWriter jsonData = ", string(jsonData))

		//封包
		dp := znet.DataPack{}
		sendMsg := znet.NewMessage(1, websocket.TextMessage, jsonData)
		encryMsg, err := dp.Pack(sendMsg)
		if err != nil {
			log.Println("client timeWriter pack err:", err, " msg:", msg)
			break
		}
		conn.WriteMessage(websocket.TextMessage, encryMsg)
		

		//发第二个消息

		i++
		if i > max {
			break
		}
	}
}
