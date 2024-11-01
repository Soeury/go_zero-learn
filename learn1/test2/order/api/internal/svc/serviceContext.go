package svc

import (
	"go_zero/learn1/test2/order/api/internal/config"
	"go_zero/learn1/test2/order/api/internal/interceptor"
	"go_zero/learn1/test2/order/model"
	"go_zero/learn1/test2/user/rpc/userinfo"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config // cache 在这里面

	UserModel model.OrderModel // mysql

	UserRPC userinfo.UserInfo
}

func NewServiceContext(c config.Config) *ServiceContext {

	// mysql 连接
	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)

	// 这里创建 rpc 客户端，可以添加拦截器
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewOrderModel(sqlxConn, c.CacheRedis), // mysql + redis cache
		UserRPC: userinfo.NewUserInfo(
			zrpc.MustNewClient(
				c.UserRPC,
				zrpc.WithUnaryClientInterceptor(interceptor.UnaryClientInterceptor), // 注意这里的函数名
			),
		),
	}
}
