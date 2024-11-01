package svc

import (
	"go_zero/learn1/test2/user/api/internal/config"
	"go_zero/learn1/test2/user/api/internal/middleware"
	"go_zero/learn1/test2/user/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

// ServiceContext 服务上下文
// rest.MiddleWare 是自定义的函数类型: func(next http.HandlerFunc) http.HandlerFunc
type ServiceContext struct {
	Config    config.Config
	Cost      rest.Middleware // 自定义中间件，注意类型，字段名字和路由中定义的中间件一致
	UserModel model.UserModel // 导入：加入 user 表
}

// NewServiceContext 这个函数返回上面的结构体实例
func NewServiceContext(c config.Config) *ServiceContext {

	//  ------- no cache -------
	//  上面的 UserModel 是一个接口类型
	//  defaultUserModel 这个结构体实现了 UserModel 接口，小写，不可跨包调用
	//  只能通过调用函数 newUserModel(conn sqlx.SqlConn) 来得到 defaultUserModel 结构体
	//  需要得到 conn sqlx.SqlConn 作为参数传入结构体

	//  ------- with cache -------
	//  创建的如果是携带缓存版本的代码
	//  首先需要修改 etc 文件和 config 文件
	//  etc 文件注意，一个 - Host 代表以台主机，这台主机的其余参数(password 等，如果有的话)不需要携带 -
	//  下面的 UserModel 的函数需要两个参数

	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:    c,
		Cost:      middleware.NewCostMiddleware().Handle,
		UserModel: model.NewUserModel(sqlxConn, c.CacheRedis),
	}
}
