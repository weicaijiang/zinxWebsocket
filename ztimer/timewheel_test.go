package ztimer

import (
	// "fmt"
	"log"
	"testing"
	"time"
)

//go test -v -run TestTimewheel
func TestTimewheel(t *testing.T) {
	//创建秒数时间轮
	secondTw := NewTimeWheel(SECOND_NAME, SECOND_INTERVAL, SECOND_SCALES, TIMERS_MAX_CAP)
	minuteTw := NewTimeWheel(MINUTE_NAME, MINUTE_INTERVAL, MINUTE_SCALES, TIMERS_MAX_CAP)
	hourTw := NewTimeWheel(HOUR_NAME, HOUR_INTERVAL, HOUR_SCALES, TIMERS_MAX_CAP)

	//将分层时间做关联
	hourTw.AddTimeWheel(minuteTw)
	minuteTw.AddTimeWheel(secondTw)

	log.Println("init timewheels done")

	//给时间轮添加定时器
	timer1 := NewTimerAfter(NewDelayFunc(myFunc, []interface{}{1, 10}), 10*time.Second)
	_ = hourTw.AddTimer(1, timer1)
	log.Println("add timer1 done")

	timer2 := NewTimerAfter(NewDelayFunc(myFunc, []interface{}{2, 20}), 20*time.Second)
	_ = hourTw.AddTimer(2, timer2)
	log.Println("add timer3 done")

	timer3 := NewTimerAfter(NewDelayFunc(myFunc, []interface{}{3, 30}), 30*time.Second)
	_ = hourTw.AddTimer(3, timer3)
	log.Println("add timer3 done")

	timer4 := NewTimerAfter(NewDelayFunc(myFunc, []interface{}{4, 40}), 40*time.Second)
	_ = hourTw.AddTimer(4, timer4)
	log.Println("add timer4 done")

	timer5 := NewTimerAfter(NewDelayFunc(myFunc, []interface{}{5, 50}), 50*time.Second)
	_ = hourTw.AddTimer(5, timer5)
	log.Println("add timer5 done")

	//时间轮运行
	secondTw.Run()
	minuteTw.Run()
	hourTw.Run()
	log.Println("timewheel are run")

	go func() {
		n := 0.0
		for {
			log.Println("tick ...", n)
			//取出最近一ms超时定时器有那些
			timers := hourTw.GetTimerWithIn(1000 * time.Microsecond)
			for _, timer := range timers {
				//调用定时方法
				timer.delayFunc.Call()
			}
			time.Sleep(500 * time.Millisecond)
			n += 0.5
		}
	}()
	time.Sleep(10 * time.Minute)
}
