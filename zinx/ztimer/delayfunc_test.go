package ztimer

import (
	"log"
	"testing"
)

func SayHello(message ...interface{}) {
	log.Println(message[0].(string), " ", message[1].(string))
}

//go test -v -run TestDelayfunc
func TestDelayfunc(t *testing.T) {
	df := NewDelayFunc(SayHello, []interface{}{"hello", "zinx websocket"})
	log.Println("df.string() = ", df.String())
	df.Call()
}
