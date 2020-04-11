package znet

import "zinxWebsocket/zinx/ziface"

type Request struct {
	//当前用户连接
	conn ziface.IConnection
	//消息封装
	message ziface.IMessage
}

//创建消息
func NewRequest(conn ziface.IConnection, msg ziface.IMessage) *Request {
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
	return r.message.GetData()
}

//得到请求类型
func (r *Request) GetMessageType() int {
	return r.message.GetMessageType()
}

//得到消息id
func (r *Request) GetMsgId() uint32 {
	return r.message.GetMsgId()
}
