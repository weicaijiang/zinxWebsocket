package main

import (
	// "encoding/base64"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"

	"strconv"

	"time"
	"zinxWebsocket/demo/message"

	"github.com/gorilla/websocket"
)

var max = 3

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

	// i := 0
	for {
		//第一个包
		// log.Println("main ReadMessage start")
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read err:", err)
			return
		}
		log.Println("main ReadMessage read server mt:", mt, " message:", string(msg[:]))

		// break
		// i++
		// if i > max {
		// 	break
		// }
	}
	//阻塞
	// select {}
	time.Sleep(60 * time.Second)
	log.Println("client exit")
}

func timeWriter(conn *websocket.Conn) {
	var i = 0
	for {
		// log.Println("WriteMessage start timeWriter i = ", i)

		//发第一个消息
		msg := &message.Account{Name: "第一个包 hello,张三", Age: i, Passwd: "123456"}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Println("client timeWriter Marshal err:", err, " msg:", msg)
			break
		}
		conn.WriteMessage(websocket.TextMessage, jsonData)

		//发第二个消息
		msg = &message.Account{Name: "第二个包 hello, 李四", Age: i, Passwd: "654321"}
		jsonData, err = json.Marshal(msg)
		if err != nil {
			log.Println("client timeWriter Marshal err:", err, " msg:", msg)
			break
		}
		conn.WriteMessage(websocket.TextMessage, jsonData)

		// //第三个是回写数据
		repeatMsg := []byte("第三个包repeat message i = " + strconv.Itoa(i))
		conn.WriteMessage(websocket.TextMessage, repeatMsg)

		//cpu阻塞下，等待读取完
		time.Sleep(5 * time.Second)

		i++
		if i > max {
			break
		}
	}

}
