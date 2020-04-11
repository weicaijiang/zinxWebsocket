package ztimer

import (
	// "fmt"
	"fmt"
	"log"
	"reflect"
)

/*
	定义一个延迟调用函数
	时间一到就调用
*/
type DelayFunc struct {
	f    func(...interface{}) //延迟函数
	args []interface{}        //延迟参数形参
}

//创建一个延迟调用函数
func NewDelayFunc(f func(...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f:    f,
		args: args,
	}
}

//打印当前回调信息
func (df *DelayFunc) String() string {
	return fmt.Sprintf("{DelayFunc: %s args: %v", reflect.TypeOf(df.f).Name(), df.args)
}

//执行调用函数
func (df *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("delayfunc call error ", err)
		}
	}()
	df.f(df.args...)
}
