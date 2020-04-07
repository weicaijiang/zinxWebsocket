package ziface

//路由接口(参考beego名字重新命名下)
type IRouter interface {
	//处理业务之前
	BeforeHandle(request IRequest)

	//处理业务
	Handle(request IRequest)

	//处理业务之后
	AfterHandle(request IRequest)
}
