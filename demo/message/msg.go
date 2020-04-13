package message

//测试结构体
type Account struct {
	Name   string `json:"name"`
	Age    int `json:"age"`
	Passwd string `json:"passwd"`
}

//测试结构体
type Room struct {
	Name   string `json:"name"`
	Port    int `json:"port"`
	Level  int `json:"level"`
}