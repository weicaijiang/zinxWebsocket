package ziface

//连接管理
type IConnManager interface {
	//添加连接
	Add(conn IConnection)

	//删除连接
	Remove(conn IConnection)

	//根据id查找连接
	Get(connID uint32) (IConnection, error)

	//总连接个数
	Len() uint32

	//清除全部连接
	ClearAll()
}
