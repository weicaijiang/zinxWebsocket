package znet

//消息
type Message struct {
	//头部 固定
	MessageType int    `json:MessageType` //消息类型 TextMessage 或 BinaryMessage之类
	Id          uint32 `json:Id`          //消息id
	//真实数据
	Data []byte `json:Data` //消息内容
}

//创建消息
func NewMessage(id uint32, mt int, data []byte) *Message {
	m := &Message{
		Id:          id,
		Data:        data,
		MessageType: mt,
	}
	return m
}

//获取消息id
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//获取消息数据
func (m *Message) GetData() []byte {
	return m.Data
}

//获取消息类型
func (m *Message) GetMessageType() int {
	return m.MessageType
}

//设置消息id
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

//设置消息数据
func (m *Message) SetData(data []byte) {
	m.Data = data
}

//设置消息类型
func (m *Message) SetMessageType(mt int) {
	m.MessageType = mt
}
