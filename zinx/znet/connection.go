package znet

import (
	"errors"
	"log"
	"net"
	"sync"
	"zinxWebsocket/zinx/utils"
	"zinxWebsocket/zinx/ziface"

	"github.com/gorilla/websocket"
)

//连接管理
type Connection struct {
	//当前属于那个server
	WsServer ziface.IServer

	//当前连接的ws
	Conn *websocket.Conn

	//连接id
	ConnID uint32

	//当前连接状态
	isClosed bool

	//告知当前连接已经退出/停止,由reder退出的信号
	ExitChan chan bool

	//无缓冲读写通信
	msgChan chan Message

	//路由管理,用来绑定msgid与api关系
	MsgHandle ziface.IMsgHandle

	//绑定属性
	property map[string]interface{}

	//保护连接属性
	propertyLock sync.RWMutex
}

//初始化连接方法
func NewConnection(server ziface.IServer, conn *websocket.Conn, connID uint32, mh ziface.IMsgHandle) *Connection {
	c := &Connection{
		WsServer:  server,
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: mh,
		isClosed:  false,
		msgChan:   make(chan Message, 1),
		ExitChan:  make(chan bool, 1),
		property:  make(map[string]interface{}),
	}

	//将当前连接放入connmgr
	c.WsServer.GetConnMgr().Add(c)

	return c
}

//读业务
func (c *Connection) StartReader() {
	log.Println("connection StartReader start connid:", c.ConnID)
	defer log.Println("connection StartReader exit connid:", c.ConnID, " remoteip:", c.Conn.RemoteAddr())
	defer c.Stop()

	//读业务
	for {
		//读取数据到内存中 messageType:TextMessage/BinaryMessage
		messageType, data, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("connection startReader read err:", err)
			break
		}
		log.Println("connection StartReader recv from client1:", string(data))

		//解包
		dp := NewDataPack()
		msg, err := dp.Unpack(messageType, data)
		if err != nil {
			log.Println("connection startReader Unpack err:", err)
			continue //下一个包可能正确
		}
		// log.Println(msg)
		log.Println("connection StartReader recv from client2:", string(msg.GetData()))

		//得到request数据
		req := &Request{
			conn:    c,
			message: msg,
		}
		//如果配置了工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandle.SendMsgToTaskQueue(req)
		} else {
			//根据gorilla/websocket官方文档 此处不能开启多线程
			// go c.MsgHandle.DoMsgHandler(req)
			c.MsgHandle.DoMsgHandler(req)
		}

	}
}

//写业务,专门发给客户端
func (c *Connection) StartWriter() {
	log.Println("connection StartWriter start")
	defer log.Println("connection StartWriter exit connid:", c.ConnID, " remoteip:", c.Conn.RemoteAddr())
	//不断的发送消息
	for {
		select {
		case msg := <-c.msgChan:
			//有数据接收
			if err := c.Conn.WriteMessage(msg.MessageType, msg.Data); err != nil {
				//写失败通知关闭连接
				log.Println("connection StartWriter err:", err)
				return
			}
		case <-c.ExitChan:
			//读出错了
			return
		}
	}
}

//启动连接，让当前连接，开始工作
func (c *Connection) Start() {
	log.Println("connection Start connid:", c.ConnID)

	//根据官方文档 读与写只能开一个线程
	//启动读数据业务
	go c.StartReader()

	//启动写数据业务
	go c.StartWriter()

	//按照开发者传递的函数来，调用回调函数
	c.WsServer.CallOnConnStart(c)
}

//停止连接，结束当前连接工作
func (c *Connection) Stop() {
	log.Println("connection stop start connid:", c.ConnID, " isClosed:", c.isClosed)
	//如是已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//按照开发者传递的函数来，调用回调函数,注意在close之前调用
	c.WsServer.CallOnConnStop(c)

	//关闭连接
	c.Conn.Close()
	//告知writer停止
	c.ExitChan <- true

	//将conn在connmgr中删除
	c.WsServer.GetConnMgr().Remove(c)

	//关闭管道
	close(c.ExitChan)
	close(c.msgChan)
	log.Println("connection stop end connid:", c.ConnID, " isClosed:", c.isClosed)
}

//获取当前连接的websocket conn
func (c *Connection) GetWsConnection() *websocket.Conn {
	return c.Conn
}

//获取当前连接的id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取连接客户端的信息，后续可以加userAgent等
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程客户端
func (c *Connection) SendMsg(messageType int, id uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection sendmsg is closed1")
	}
	dp := NewDataPack()
	message := NewMessage(id, messageType, data)
	decryMsg, err := dp.Pack(message)
	if err != nil {
		return errors.New("connection sendmsg Pack error")
	}

	//重新设置下加密后的data
	message.SetData(decryMsg)
	log.Println("connection sendmsg isclosed = ", c.isClosed)

	//发消息给通道
	c.msgChan <- *message

	return nil
}

//设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

//获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("connection getproperty get error key:" + key)
	}
}

//移除设置属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
