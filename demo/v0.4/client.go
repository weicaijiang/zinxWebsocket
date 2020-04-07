package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	// "time"

	"github.com/gorilla/websocket"
)

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
		log.Println("ReadMessage start")
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("read err:", err)
			return
		}
		log.Println("read server message", string(message[:]))
		// break
		i++
		if i > 300 {
			break
		}
	}

}

func timeWriter(conn *websocket.Conn) {
	var i = 0
	for {
		log.Println("WriteMessage start timeWriter i = ", i)
		// time.Sleep(time.Second * 1)
		// conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
		conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(i)))
		i++
		if i > 300 {
			break
		}
	}
}
