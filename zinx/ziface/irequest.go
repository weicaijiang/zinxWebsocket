package ziface

//把客户端数据包装成一个requst
type IRequest interface {
	//得到当前连接
	GetConnection() IConnection
	//得到请求数据
	GetData() string
}
