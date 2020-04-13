package znet

import "zinxWebsocket/zinx/ziface"

type Request struct {
	//当前用户连接
	conn ziface.IConnection
	//消息封装
	message string
}

//创建消息
func NewRequest(conn ziface.IConnection, msg string) *Request {
	r := &Request{
		conn:    conn,
		message: msg,
	}
	return r
}

//得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

//得到请求数据
func (r *Request) GetData() string {
	return r.message
}
