# zinxwebsocket

根据zinx写的go语言版 websocket,官方没websocket版本，这里模仿一个
1.去掉了imessage,idatapack,把message解放出来
2.服务器路由默认是BaseHandler
3.日志输出 统一是 文件名 + 函数名 + 信息

依赖
go get -u github.com/gorilla/websocket

运行
//服务器
cd demo  
go run .\server.go
客户端
cd demo 
go run .\client.go

gitee地址 https://gitee.com/sundayme/zinxWebsocket
github地址 https://github.com/weicaijiang/zinxWebsocket

tcp版本请查看  https://github.com/aceld/zinx.git