package zlog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//日志全部方法

const (
	LOG_MAX_BUFF = 1024 * 1024
)

//日志头部信息标记
const (
	BitDate         = 1 << iota                            //日期标记位 2020/04/10
	BitTime                                                //时间标记位 09:13:12
	BitMicroSeconds                                        //微秒级标记位 09:13:12.258258
	BitLongFile                                            //完整文件名 /home/go/src/zinx/server.log
	BitShortFile                                           //最后文件名 server.log
	BitLevel                                               //当前日志级别 0(Debug) 1(Info) 2(Warn) 3(Error) 4(Panic) 5(Fatal)
	BitStdFlag      = BitDate | BitTime                    //标准头部日志格式
	BitDefault      = BitLevel | BitShortFile | BitStdFlag //默认日志头部格式
)

//日志级别
const (
	LogDebug = iota
	LogInfo
	LogWarn
	LogError
	LogPanic
	LogFatal
)

//日志级别对应的显示字符
var levels = []string{
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

type ZinxLogger struct {
	//确保多线程读写
	mu sync.Mutex
	//每行log日志前缀字符串
	prefix string
	//日志标记
	flag int
	//日志输出的文件描述
	out io.Writer
	//输出缓冲区
	buf bytes.Buffer
	//当前日志绑定输出文件
	file *os.File
	//是否打印调试信息
	debugClose bool
	//获取日志文件和代码上述的runtime.call函数调用信息
	callDepth int
}

/*
	创建一个日志
	out 标准输出的io
	prefix 日志的前缀
	flag 当前日志头部信息
*/
func NewZinxLog(out io.Writer, prefix string, flag int) *ZinxLogger {
	//默认打开debug calldepth 深度为2层
	zlog := &ZinxLogger{out: out, prefix: prefix, flag: flag, file: nil, debugClose: false, callDepth: 2}
	//设置log对象回收资源 （不设置也可以 强迫症来了没办法）
	runtime.SetFinalizer(zlog, CleanZinxLog)
	return zlog
}

//日志回收处理
func CleanZinxLog(log *ZinxLogger) {
	log.closeFile()
}

//制作当前日志数据格式信息
func (log *ZinxLogger) formatHeader(buf *bytes.Buffer, t time.Time, file string, line int, level int) {
	//如果当前前缀不为空，那么需要行写前缀
	if log.prefix != "" {
		buf.WriteByte('<')
		buf.WriteString(log.prefix)
		buf.WriteByte('>')
	}
	//已经设置了时间相关的标识位，那么需要添加时间信息日志头部
	if log.flag&(BitDate|BitTime|BitMicroSeconds) != 0 {
		//日期标记位
		if log.flag&BitDate != 0 {
			year, month, day := t.Date()
			// fmt.Println("year:", year, " month:", month, " day:", day)
			intToWidth(buf, year, 4)
			buf.WriteByte('/') //2019/
			intToWidth(buf, int(month), 2)
			buf.WriteByte('/') //2019/04
			intToWidth(buf, day, 2)
			buf.WriteByte(' ') //2019/04/08
		}

		//时间位被标记
		if log.flag&(BitTime|BitMicroSeconds) != 0 {
			hour, min, sec := t.Clock()
			intToWidth(buf, hour, 2)
			buf.WriteByte(':')
			intToWidth(buf, min, 2)
			buf.WriteByte(':')
			intToWidth(buf, sec, 2)
			if log.flag&BitMicroSeconds != 0 {
				buf.WriteByte('.')
				intToWidth(buf, t.Nanosecond()/1e3, 6)
			}
			buf.WriteByte(' ')
		}

		//日志级别
		if log.flag&BitLevel != 0 {
			buf.WriteString(levels[level])
		}
		//日志当前代码调用文件名名称标记
		if log.flag&(BitLongFile|BitShortFile) != 0 {
			//短文件名
			if log.flag&BitShortFile != 0 {
				short := file
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						//找到最后一个 / 之后的名字
						short = file[i+1:]
						break
					}
				}
				file = short
			}
			buf.WriteString(file)
			buf.WriteByte(':')
			intToWidth(buf, line, -1) //行数
			buf.WriteString(": ")
		}
	}
}

//输出日志文件 原方法
func (log *ZinxLogger) OutPut(level int, s string) error {
	now := time.Now()
	var file string //文件名
	var line int    // 当前代码行
	log.mu.Lock()
	defer log.mu.Unlock()

	if log.flag&(BitShortFile|BitLongFile) != 0 {
		log.mu.Unlock()
		var ok bool
		//得到当前调用者的文件名称和执行代码行数
		_, file, line, ok = runtime.Caller(log.callDepth)
		if !ok {
			file = "unknownfile"
			line = 0
		}
		log.mu.Lock()
	}
	//清零buf
	log.buf.Reset()
	//写日志头
	log.formatHeader(&log.buf, now, file, line, level)
	//写日志内容
	log.buf.WriteString(s)
	//补充回车
	if len(s) > 0 && s[len(s)-1] != '\n' {
		log.buf.WriteByte('\n')
	}
	//将填充好的内容输出到io上
	_, err := log.out.Write(log.buf.Bytes())
	return err
}

//debug
func (log *ZinxLogger) Debugf(format string, v ...interface{}) {
	if log.debugClose {
		return
	}
	_ = log.OutPut(LogDebug, fmt.Sprintf(format, v...))
}

func (log *ZinxLogger) Debug(v ...interface{}) {
	if log.debugClose {
		return
	}
	_ = log.OutPut(LogDebug, fmt.Sprintln(v...))
}

//info
func (log *ZinxLogger) Infof(format string, v ...interface{}) {
	_ = log.OutPut(LogInfo, fmt.Sprintf(format, v...))
}

func (log *ZinxLogger) Info(v ...interface{}) {
	_ = log.OutPut(LogInfo, fmt.Sprintln(v...))
}

//Warn
func (log *ZinxLogger) Warnf(format string, v ...interface{}) {
	_ = log.OutPut(LogWarn, fmt.Sprintf(format, v...))
}

func (log *ZinxLogger) Warn(v ...interface{}) {
	_ = log.OutPut(LogWarn, fmt.Sprintln(v...))
}

//Error
func (log *ZinxLogger) Errorf(format string, v ...interface{}) {
	_ = log.OutPut(LogError, fmt.Sprintf(format, v...))
}

func (log *ZinxLogger) Error(v ...interface{}) {
	_ = log.OutPut(LogError, fmt.Sprintln(v...))
}

//Panic
func (log *ZinxLogger) Panicf(format string, v ...interface{}) {
	_ = log.OutPut(LogPanic, fmt.Sprintf(format, v...))
}

func (log *ZinxLogger) Panic(v ...interface{}) {
	_ = log.OutPut(LogPanic, fmt.Sprintln(v...))
}

//Fatal
func (log *ZinxLogger) Fatalf(format string, v ...interface{}) {
	_ = log.OutPut(LogFatal, fmt.Sprintf(format, v...))
}

func (log *ZinxLogger) Fatal(v ...interface{}) {
	_ = log.OutPut(LogFatal, fmt.Sprintln(v...))
}

//stack
func (log *ZinxLogger) Stack(v ...interface{}) {
	s := fmt.Sprint(v...)
	s += "\n"
	buf := make([]byte, LOG_MAX_BUFF)
	n := runtime.Stack(buf, true) //得到当前堆栈信息
	s += string(buf[:n])
	s += "\n"
	_ = log.OutPut(LogError, s)
}

//获取当前日志标记
func (log *ZinxLogger) Flags() int {
	log.mu.Lock()
	defer log.mu.Unlock()
	return log.flag
}

//重置当前日志标记
func (log *ZinxLogger) ResetFlags(flag int) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.flag = flag
}

//添加flag标记
func (log *ZinxLogger) AddFlags(flag int) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.flag |= flag
}

//设置日志前缀
func (log *ZinxLogger) SetPrefix(prefix string) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.prefix = prefix
}

//设置日志文件输出
func (log *ZinxLogger) SetLogFile(fileDir string, fileName string) {
	var file *os.File
	//创建日志文件夹
	_ = mkdirLog(fileDir)

	fullPath := fileDir + "/" + fileName
	if log.checkFileExist(fullPath) {
		//文件存在
		file, _ = os.OpenFile(fullPath, os.O_APPEND|os.O_RDWR, 0644)
	} else {
		//文件不存在
		file, _ = os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	}
	log.mu.Lock()
	defer log.mu.Unlock()
	//关闭之前绑定的文件
	log.closeFile()
	log.file = file
	log.out = file
}

//关闭日志绑定文件
func (log *ZinxLogger) closeFile() {
	if log.file != nil {
		_ = log.file.Close()
		log.file = nil
		log.out = os.Stderr
	}
}

//debug open/close
func (log *ZinxLogger) CloseDebug() {
	log.debugClose = true
}

func (log *ZinxLogger) OpenDebug() {
	log.debugClose = false
}

//一些工具方法
//判断日志是否存在
func (log *ZinxLogger) checkFileExist(filename string) bool {
	exist := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//创建目录
func mkdirLog(dir string) error {
	_, err := os.Stat(dir)
	b := err == nil || os.IsExist(err)
	if !b {
		if err = os.MkdirAll(dir, 0775); err != nil {
			if os.IsPermission(err) {
				// err = err
			}
		}
	}
	return err
}

//将一个整形转成一个固定长度的字符
//要确保buffer有容量空间
func intToWidth(buf *bytes.Buffer, i int, wid int) {
	//简单就好，不需要复杂的格式
	buf.WriteString(strconv.Itoa(i))

	// var u uint = uint(1)
	// if u == 0 && wid <= 1 {
	// 	buf.WriteByte('0')
	// 	return
	// }

	// var b [32]byte
	// bp := len(b)
	// for ; u > 0 || wid > 0; u /= 10 {
	// 	bp--
	// 	wid--
	// 	b[bp] = byte(u%10) + '0'
	// }

	// for bp < len(b) {
	// 	buf.WriteByte(b[bp])
	// 	bp++
	// }
}
