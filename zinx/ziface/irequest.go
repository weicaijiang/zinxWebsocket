package ziface

//把客户端数据包装成一个requst
type IRequest interface {
	//得到当前连接
	GetConnection() IConnection

	//得到请求数据
	GetData() []byte

	//得到请求类型
	GetMessageType() int

	//得到消息id
	GetMsgId() uint32
}
