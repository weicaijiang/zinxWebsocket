package zlog

import (
	// "go/format"
	"os"
)

// 全局提供一个log对外句柄 直接用api使用

var StdZinxLog = NewZinxLog(os.Stderr, "", BitDefault)

//获取 标记位
func Flags() int {
	return StdZinxLog.Flags()
}

//设置 标记位
func ResetFlags(flag int) {
	StdZinxLog.ResetFlags(flag)
}

//添加 标记位
func AddFlags(flag int) {
	StdZinxLog.AddFlags(flag)
}

//设置 前缀
func SetPrefix(prefix string) {
	StdZinxLog.SetPrefix(prefix)
}

//设置绑定日志文件
func SetLogFile(fileDir string, fileName string) {
	StdZinxLog.SetLogFile(fileDir, fileName)
}

//设置 关闭debug
func CloseDebug() {
	StdZinxLog.CloseDebug()
}

func OpenDebug() {
	StdZinxLog.OpenDebug()
}

//Debug
func Debugf(format string, v ...interface{}) {
	StdZinxLog.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	StdZinxLog.Debug(v...)
}

//Warn
func Warnf(format string, v ...interface{}) {
	StdZinxLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	StdZinxLog.Warn(v...)
}

//Info
func Infof(format string, v ...interface{}) {
	StdZinxLog.Infof(format, v...)
}

func Info(v ...interface{}) {
	StdZinxLog.Info(v...)
}

//Error
func Errorf(format string, v ...interface{}) {
	StdZinxLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	StdZinxLog.Error(v...)
}

//Panic
func Panicf(format string, v ...interface{}) {
	StdZinxLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	StdZinxLog.Panic(v...)
}

//Fatal
func Fatalf(format string, v ...interface{}) {
	StdZinxLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	StdZinxLog.Fatal(v...)
}

//stack
func Stack(v ...interface{}) {
	StdZinxLog.Stack(v...)
}

func init() {
	//因为stdzinxlog多调用了一层，这边多一层
	StdZinxLog.callDepth = 3
}
