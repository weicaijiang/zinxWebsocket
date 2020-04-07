package ziface

//封包，拆包模块
//数据必须是 {id:200,data:""} 之后再base64WriteMessage编码
type IDataPack interface {
	//封包
	Pack(msg IMessage) ([]byte, error)
	//拆包
	Unpack(messageType int, data []byte) (IMessage, error)
}
