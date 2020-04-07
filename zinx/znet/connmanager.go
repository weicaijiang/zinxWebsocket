package znet

import (
	"errors"
	"log"
	"sync"
	"zinxWebsocket/zinx/ziface"
)

//连接管理
type ConnManager struct {
	//管理连接
	connnections map[uint32]ziface.IConnection
	//保护连接锁
	connLock sync.RWMutex
}

//创建管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connnections: make(map[uint32]ziface.IConnection),
	}
}

//添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//共享锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connMgr.connnections[conn.GetConnID()] = conn
	log.Println("connmanager add connid:", conn.GetConnID())
}

//删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//共享锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connnections, conn.GetConnID())
	log.Println("connmanager Remove connid:", conn.GetConnID())
}

//根据id查找连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//共享锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connnections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connmanager Get err")
	}
}

//总连接个数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connnections)
}

//清除全部连接
func (connMgr *ConnManager) ClearConn() {
	//共享锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	//删除conn停止工作
	for connID, conn := range connMgr.connnections {
		conn.Stop()
		delete(connMgr.connnections, connID)
	}
	log.Println("connmanager ClearConn success")
}
