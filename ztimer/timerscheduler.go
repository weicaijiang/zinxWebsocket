package ztimer

import (
	"log"
	"math"
	"sync"
	"time"
)

//时间轮调度器
//依赖模块 delayfunc.go timer.go timewheel.go
const (
	//默认缓冲触发的队列大小
	MAX_CHAN_BUF = 2048
	//默认最大误差时间
	MAX_TIME_DELAY = 100
)

type TimerScheduler struct {
	//当前调度器的最高优先级
	tw *TimeWheel
	//定时器编号累加器
	idGen uint32
	//已经触发定时器channel
	triggerChan chan *DelayFunc
	//互斥锁
	sync.RWMutex
}

//返回一个定时调度器 主要创建分层定时器 并做关联 依次启动
func NewTimerScheduler() *TimerScheduler {
	//创建秒数时间轮
	secondTw := NewTimeWheel(SECOND_NAME, SECOND_INTERVAL, SECOND_SCALES, TIMERS_MAX_CAP)
	minuteTw := NewTimeWheel(MINUTE_NAME, MINUTE_INTERVAL, MINUTE_SCALES, TIMERS_MAX_CAP)
	hourTw := NewTimeWheel(HOUR_NAME, HOUR_INTERVAL, HOUR_SCALES, TIMERS_MAX_CAP)

	//将分层时间做关联
	hourTw.AddTimeWheel(minuteTw)
	minuteTw.AddTimeWheel(secondTw)

	//时间轮运行
	secondTw.Run()
	minuteTw.Run()
	hourTw.Run()
	log.Println("timerscheduler NewTimerScheduler are run")

	return &TimerScheduler{
		tw:          hourTw,
		triggerChan: make(chan *DelayFunc, MAX_CHAN_BUF),
	}
}

//创建一个定点timer 并将timer添加到分层时间轮中 返回 timer中的tid
func (ts *TimerScheduler) CreateTimerAt(df *DelayFunc, unixNano int64) (uint32, error) {
	ts.Lock()
	defer ts.Unlock()
	ts.idGen++
	return ts.idGen, ts.tw.AddTimer(ts.idGen, NewTimerAt(df, unixNano))
}

//创建一个定时器，并将timer添加到分层时间轮中，返回timer的tid
func (ts *TimerScheduler) CreateTimerAfter(df *DelayFunc, duration time.Duration) (uint32, error) {
	ts.Lock()
	defer ts.Unlock()
	ts.idGen++
	return ts.idGen, ts.tw.AddTimer(ts.idGen, NewTimerAfter(df, duration))
}

//删除timer
func (ts *TimerScheduler) CancelTimer(tid uint32) {
	ts.Lock()
	ts.Unlock()
	ts.tw.RemoveTimer(tid)
}

//获取时间结束的延迟执行函数
func (ts *TimerScheduler) GetTriggerChan() chan *DelayFunc {
	return ts.triggerChan
}

//非阻塞方式启动timerscheduler
func (ts *TimerScheduler) Start() {
	go func() {
		for {
			//当前时间
			now := UinxMilli()
			//获取最近的定时器集合
			timerList := ts.tw.GetTimerWithIn(MAX_TIME_DELAY * time.Millisecond)
			for _, timer := range timerList {
				if math.Abs(float64(now-timer.unixts)) > MAX_TIME_DELAY {
					//已经超时，报警处理
					log.Println("timerscheduler start call at ", timer.unixts, " real call at ", now, " delay ", now-timer.unixts)
				}
				ts.triggerChan <- timer.delayFunc
			}
			time.Sleep(MAX_TIME_DELAY / 2 * time.Millisecond)
		}
	}()
}

//时间轮定时器 自动调度
func NewAutoExecTimerScheduler() *TimerScheduler {
	//创建一个调度器
	autoExecScheduler := NewTimerScheduler()
	//启动调度器
	autoExecScheduler.Start()

	//阻塞获取定时器并执行
	go func() {
		delayFunc := autoExecScheduler.GetTriggerChan()
		for df := range delayFunc {
			go df.Call()
		}
	}()
	return autoExecScheduler
}
