# zinx-websocket

根据zinx写的websocket,官方没websocket版本，这里模仿一个
比如发 msgid = 1 data = "张三"
内部分封包成 {Id:1,MessageType:1,Data:"张三"}
头部固定会添加 Id 与 MessageType两个字段

日志输出 统一是 文件名 + 函数名 + 信息
