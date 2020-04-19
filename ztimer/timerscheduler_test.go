package ztimer

import (
	"log"
	"testing"
	"time"
)

//触发函数
func foo(args ...interface{}) {
	log.Println("i am no ", args[0].(int), " function delay ", args[1].(int))
}

//手动创建调度运行时间轮 go test -v -run TestNewTimerScheduler
func TestNewTimerScheduler(t *testing.T) {
	timerScheduler := NewTimerScheduler()
	timerScheduler.Start()

	//在scheduler中添加timer
	for i := 1; i < 2000; i++ {
		f := NewDelayFunc(foo, []interface{}{i, i * 3})
		tid, err := timerScheduler.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			log.Println("create timer error ", tid, err)
			break
		}
	}
	//执行调度器
	go func() {
		delayFuncChan := timerScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			df.Call()
		}
	}()
	select {}
}

//自动调度时间轮 go test -v -run TestNewAutoExecTimerScheduler
func TestNewAutoExecTimerScheduler(t *testing.T) {
	autoTs := NewAutoExecTimerScheduler()

	//给调度器添加timer
	for i := 0; i < 2000; i++ {
		f := NewDelayFunc(foo, []interface{}{i, i * 3})
		tid, err := autoTs.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			log.Println("create timer error ", tid, err)
			break
		}

	}
	select {}
}
