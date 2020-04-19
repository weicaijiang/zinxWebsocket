package ztimer

import "time"

const (
	HOUR_NAME     = "HOUR"
	HOUR_INTERVAL = 60 * 60 * 1e3 //ms为精度
	HOUR_SCALES   = 12

	MINUTE_NAME     = "MINUTE"
	MINUTE_INTERVAL = 60 * 1e3
	MINUTE_SCALES   = 60

	SECOND_NAME     = "SECOND"
	SECOND_INTERVAL = 1e3
	SECOND_SCALES   = 60

	TIMERS_MAX_CAP = 2048 //每个时间内运行定时器的最大数
)

/*
 有关时间的几个换算
 time.Second(秒) = time.Millisecond*1e3
 time.Millisecond(毫秒) = time.Microsecond*1e3
 time.Microsecond(微秒) = time.Nanosecond*1e3

 time.Now().UnixNano() ==> time.Nanosecond(纳秒)
*/

//定时器实现
type Timer struct {
	//延迟调用函数
	delayFunc *DelayFunc
	//调用时间(unix时间 单位ms毫秒)
	unixts int64
}

//返回1970-1-1到今天的毫秒数
func UinxMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

/*创建一个定时器，在指定的时间调用
unixNano:1970-1-1到今天的纳秒数
*/
func NewTimerAt(df *DelayFunc, unixNano int64) *Timer {
	return &Timer{
		delayFunc: df,
		unixts:    unixNano / 1e6,
	}
}

//创建一个定时器，在当前延迟duration纳秒之后触发
func NewTimerAfter(df *DelayFunc, duration time.Duration) *Timer {
	return NewTimerAt(df, time.Now().UnixNano()+int64(duration))
}

//启动定时器，用一个go承载
func (t *Timer) Run() {
	go func() {
		now := UinxMilli()
		//设定的时间在当前时间之后
		if t.unixts > now {
			//睡眠，直到时间超时，以微秒为单位进行睡眠
			time.Sleep(time.Duration(t.unixts-now) * time.Millisecond)
		}
		//调用之前的注册好的函数方法
		t.delayFunc.Call()
	}()
}
