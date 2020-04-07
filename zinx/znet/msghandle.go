package znet

import (
	"log"
	"strconv"
	"zinxWebsocket/zinx/utils"
	"zinxWebsocket/zinx/ziface"
)

//消息处理模块
type MsgHandle struct {
	//每一个msgid对应处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker取任务消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

//创建
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置获取
	}
}

//执行对应的路由处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		log.Println("msghandle DoMsgHandler not register msgid:", request.GetMsgId())
		return
	}
	handler.BeforeHandle(request)
	handler.Handle(request)
	handler.AfterHandle(request)
}

//添加路由
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//是否存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("msghandle addrouter repeate error msgid: " + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
}

//册除路由
func (mh *MsgHandle) RemoveRouter(msgId uint32) {
	_, ok := mh.Apis[msgId]
	if !ok {
		log.Println("msghandle RemoveRouter not register msgid:", msgId)
		return
	}
	log.Println("msghandle RemoveRouter success msgid:", msgId)
	delete(mh.Apis, msgId)
}

//启动一个worker工作池 只能发生一次
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerpoolsize启动一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//当前worker对应的chan开启空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前的worker，阻塞等待
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	log.Println("msghandle StartOneWorker workerid:", workerID)
	for {
		select {
		//如果有消息过来，执行绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给taskqueue处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配
	workerID := request.GetConnection().GetConnID() % utils.GlobalObject.WorkerPoolSize
	log.Println("msghandle SendMsgToTaskQueue workerID:", workerID, " connid:", request.GetConnection().GetConnID())
	//将消息发送给对应的worker
	mh.TaskQueue[workerID] <- request
}
