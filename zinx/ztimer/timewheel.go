package ztimer

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

/*
	一个网络服务器程序时需要管理大量客户连接
	其中每个客户端连接都需要管理它的timeout
	通常连接超时一般设置为30-60秒，并不需要太精确的时间控制
	另外由于服务器管理着多达数万个不现的连接
	我们没法为每个连接使用一个timer，太橇资源

	用时间轮的方式来管理和维护大量的timer即可
*/
type TimeWheel struct {
	//名称
	name string
	//刻度时间间隔 单位ms
	interval int64
	//每个时间轮上的刻度数
	scales int
	//当前时间指针指向
	curIndex int
	//每个刻度所存放的timer最大容量
	maxCap int
	//当前时间上所有timer int表示当前时间轮数 uint32表示timer的id号
	timerQueue map[int]map[uint32]*Timer
	//下一层时间轮
	nextTimeWheel *TimeWheel
	//互斥锁 继承此方法
	sync.RWMutex
}

/*
	创建一个时间轮
	name 时间轮的名字
	interval 每个刻度之间的durationg时间间隔
	scales 当前时间轮有多少个刻度 比如一个小时12个刻度
	maxcap 每个刻度保存的最大timer个数
*/
func NewTimeWheel(name string, interval int64, scales int, maxCap int) *TimeWheel {
	tw := &TimeWheel{
		name:       name,
		interval:   interval,
		scales:     scales,
		maxCap:     maxCap,
		timerQueue: make(map[int]map[uint32]*Timer, scales),
	}
	//初始化map
	for i := 0; i < scales; i++ {
		tw.timerQueue[i] = make(map[uint32]*Timer, maxCap)
	}
	log.Println("timewheel NewTimeWheel init success name:", tw.name)
	return tw
}

/*
	将一个timer定时器加入到分层时间轮中
	tid 每个定时器的唯一标识
	t 当前被加入的时间轮定时器
	forceNext 是否强制将定时器添加到下一层时间轮中

	算法：如果当前timer超时时间大于一个刻度 那到进行hash计算 找到对应的刻度上添加
*/
func (tw *TimeWheel) addTimer(tid uint32, t *Timer, forceNext bool) error {
	defer func() error {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("timewheel addtimer err: %s", err)
			log.Println(errStr)
			return errors.New(errStr)
		}
		return nil
	}()

	//得到当前超时时间间隔，单位为毫秒
	delayInterval := t.unixts - UinxMilli()

	//如果当前的超时时间大于一个刻度时间
	if delayInterval >= tw.interval {
		//得到需要跨越几个刻度
		dn := delayInterval / tw.interval
		//在对应的刻度上定时器timer加入当前定时器，因为是环形要求余
		tw.timerQueue[(tw.curIndex + int(dn)%tw.scales)][tid] = t
		return nil
	}
	//如是当前超时时间小于一个刻度，并且当前时间没下一层，经度最小的时间轮
	if delayInterval < tw.interval && tw.nextTimeWheel == nil {
		if forceNext {
			/*
				如果设置为强制移到下一刻度，那么将定时器移到下一刻度
				这种情况主要是时间自动轮转的情况
				这时底层时间轮，该定时器在转动时，如果没有取走，那么该定时器不会再被发现
				时间轮已经过去，如果不强制把该定时器timer移到下一时刻，就永远不会触发调用
				这里强制将timer移到下个刻度集合中，等待调用者下次轮转之前取走该定时器
			*/
			tw.timerQueue[(tw.curIndex+1)%tw.scales][tid] = t
		} else {
			//如果手动添加定时器，那么直接将timer添加到对应底层时间轮的当前集合中
			tw.timerQueue[tw.curIndex][tid] = t
		}
		return nil
	}
	//如果当前超时时间小于一个刻度间隔并且有下一轮时间轮
	if delayInterval < tw.interval {
		return tw.nextTimeWheel.AddTimer(tid, t)
	}
	return nil
}

//添加一个timer到一个时间轮中（非时间轮自转情况）
func (tw *TimeWheel) AddTimer(tid uint32, t *Timer) error {
	tw.Lock()
	defer tw.Unlock()
	return tw.addTimer(tid, t, false)
}

//删除一个定时器
func (tw *TimeWheel) RemoveTimer(tid uint32) {
	tw.Lock()
	defer tw.Unlock()

	for i := 0; i < tw.scales; i++ {
		if _, ok := tw.timerQueue[i][tid]; ok {
			delete(tw.timerQueue[i], tid)
			//有可能在多轮中有此tid，此处不能break
		}
	}
}

//给一个时间轮添加下层时间轮，比如给小时时间轮添加分钟时间轮，给分钟添加秒时间轮
func (tw *TimeWheel) AddTimeWheel(next *TimeWheel) {
	tw.nextTimeWheel = next
	log.Println("timewheel addtimewheel nowname:", tw.name, " nextname:", next.name)
}

//启动时间轮
func (tw *TimeWheel) run() {
	for {
		//时间轮每间隔interval一刻度时间，触发一次转动
		time.Sleep(time.Duration(tw.interval) * time.Millisecond)

		tw.Lock()
		//取出挂载在当前刻度的全部定时器
		curTimers := tw.timerQueue[tw.curIndex]
		//当前定时器要重新添加，给当前刻度再重新开辟一个map timer 容器
		tw.timerQueue[tw.curIndex] = make(map[uint32]*Timer, tw.maxCap)
		for tid, timer := range curTimers {
			//这里属于时间轮自动转动，forceNext设置为true
			tw.addTimer(tid, timer, true)
		}
		//取出下一个刻度挂载全部定时器 进行重新加载
		nextTimers := tw.timerQueue[(tw.curIndex+1)%tw.scales]
		tw.timerQueue[(tw.curIndex+1)%tw.scales] = make(map[uint32]*Timer, tw.maxCap)
		for tid, timer := range nextTimers {
			tw.addTimer(tid, timer, true)
		}
		//当前刻度走下一轮
		tw.curIndex = (tw.curIndex + 1) % tw.scales
		tw.Unlock()
	}
}

//非阻塞方式让时间轮转起来
func (tw *TimeWheel) Run() {
	go tw.run()
	log.Println("timewheel Run name:", tw.name, " is running")
}

//获取在一定时间内的timer
func (tw *TimeWheel) GetTimerWithIn(duration time.Duration) map[uint32]*Timer {
	//最终触发的定时器一定挂载在最低层的定时器上
	//1.找到最低层的时间轮
	leaftw := tw
	for leaftw.nextTimeWheel != nil {
		leaftw = leaftw.nextTimeWheel
	}

	leaftw.Lock()
	defer leaftw.Unlock()
	//返回timer集合
	timerList := make(map[uint32]*Timer)

	now := UinxMilli()
	//取出当前时间轮的全部timer
	for tid, timer := range leaftw.timerQueue[leaftw.curIndex] {
		if timer.unixts-now < int64(duration/time.Millisecond) {
			//当前定时器已经超时
			timerList[tid] = timer
			//定时器已被超时取走，从当前时间轮上移除此定时器
			delete(leaftw.timerQueue[leaftw.curIndex], tid)
		}
	}
	return timerList
}
