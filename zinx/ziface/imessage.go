package ziface

//消息定义
type IMessage interface {
	//获取消息id
	GetMsgId() uint32
	//获取消息数据
	GetData() []byte
	//获取消息类型
	GetMessageType() int

	//设置消息id
	SetMsgId(id uint32)
	//设置消息数据
	SetData(data []byte)
	//设置消息类型
	SetMessageType(mt int)
}
