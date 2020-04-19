package zlog

import "testing"

func TestStdZlog(t *testing.T) {
	Debug("zinx debug content1")
	Infof("zinx info a = %d\n", 10)

	ResetFlags(BitDate | BitLongFile | BitLevel)
	Info("zinx info 2")

	SetPrefix("Module")
	Error("zinx error content")

	AddFlags(BitShortFile | BitTime | BitMicroSeconds)
	// Stack("zinx stack")

	SetLogFile("./log", "testfile.log")
	Debug("zinx debug content666")
	Error("zinx debug content666")

	CloseDebug()
	Debug("i am not in")
	Error("zinx error after debug")
}
