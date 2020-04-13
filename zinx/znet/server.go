package znet

import (
	// "errors"
	"log"
	"net/http"
	"strconv"
	"zinxWebsocket/zinx/utils"
	"zinxWebsocket/zinx/ziface"

	"github.com/gorilla/websocket"
)

//服务器实现 ws://127.0.0.1:8080/echo
type Server struct {
	//服务器名称
	Name string
	//服务器协议 ws,wss
	Scheme string
	//服务器ip地址
	Host string
	//服务器端口
	Port uint32
	//协议
	Path string
	//路由管理,用来绑定msgid与api关系
	MsgHandle ziface.IMsgHandle
	//连接属性
	ConnMgr ziface.IConnManager
	//连接回调
	OnConnStart func(ziface.IConnection)
	//关闭回调
	OnConnStop func(ziface.IConnection)
}

//连接信息
var upgrader = websocket.Upgrader{
	ReadBufferSize:  int(utils.GlobalObject.MaxPackageSize), //读取最大值
	WriteBufferSize: int(utils.GlobalObject.MaxPackageSize), //写最大值
	//解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//全局conectionid 后续使用uuid生成
var cid uint32

//websocket回调
func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("server wsHandler upgrade err:", err)
		return
	}
	// defer log.Println("server wsHandler client is closed")
	// defer conn.Close()

	// 判断是否超出个数
	if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
		//todo 给用户发一个关闭连接消息
		log.Println("server wsHandler too many connection")
		conn.Close()
		return
	}

	log.Println("server wsHandler a new client coming ip:", conn.RemoteAddr())
	//处理新连接业务方法
	dealConn := NewConnection(s, conn, cid, s.MsgHandle)
	go dealConn.Start()
	cid++
}

//启动
func (s *Server) Start() {
	log.Println("server start name:", utils.GlobalObject.Name, " scheme:", s.Scheme, " ip:", s.Host, " port:", strconv.Itoa(int(s.Port)), " path:", s.Path,
		" MaxConn:", utils.GlobalObject.MaxConn, " MaxPackageSize:", utils.GlobalObject.MaxPackageSize,
		" WorkerPoolSize:",utils.GlobalObject.WorkerPoolSize)
	//开启工作线程
	s.MsgHandle.StartWorkerPool()
	
	http.HandleFunc("/"+s.Path, s.wsHandler)
	err := http.ListenAndServe(s.Host+":"+strconv.Itoa(int(s.Port)), nil)
	if err != nil {
		log.Println("server start listen error:", err)
	}
}

//停止
func (s *Server) Stop() {
	log.Println("server stop name:", s.Name)
	//停止所有连接
	s.ConnMgr.ClearConn()
}

//运行状态
func (s *Server) Serve() {
	s.Start()

	//额外的工作

}

//添加路由
func (s *Server) SetRouter( router ziface.IRouter) {
	s.MsgHandle.SetRouter(router)
}

//返回 连接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//连接之前回调
func (s *Server) SetOnConnStart(hookStart func(conn ziface.IConnection)) {
	s.OnConnStart = hookStart
}

//关闭之前回调
func (s *Server) SetOnConnStop(hookStop func(conn ziface.IConnection)) {
	s.OnConnStop = hookStop
}

//调用连接之前
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart == nil {
		log.Println("server CallOnConnStart error is nil")
		return
	}
	s.OnConnStart(conn)
}

//调用关闭之前
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart == nil {
		log.Println("server CallOnConnStop error is nil")
		return
	}
	s.OnConnStop(conn)
}

//初始化
func NewServer() ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		Scheme:    utils.GlobalObject.Scheme,
		Host:      utils.GlobalObject.Host,
		Port:      utils.GlobalObject.Port,
		Path:      utils.GlobalObject.Path, // 比如 /echo
		MsgHandle: NewMsgHandle(),
		ConnMgr:   NewConnManager(),
	}
	return s
}
