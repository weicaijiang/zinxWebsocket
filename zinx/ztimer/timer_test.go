package ztimer

import (
	"log"
	"testing"
	"time"
)

//定义一个超时函数
func myFunc(v ...interface{}) {
	log.Println("No.", v[0].(int), " function called delay ", v[1].(int), " second")
}

//go test -v -run TestTimer
func TestTimer(t *testing.T) {
	for i := 0; i < 5; i++ {
		go func(i int) {
			NewTimerAfter(NewDelayFunc(myFunc, []interface{}{i, 2 * i}), time.Duration(2*i)*time.Second).Run()
		}(i)
	}
	//主进程等待其它go，由于run方法是用另一个go承载，这里不能使用waitGroup
	time.Sleep(1 * time.Minute)
}
