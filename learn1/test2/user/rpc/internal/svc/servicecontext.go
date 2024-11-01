package svc

import (
	"go_zero/learn1/test2/user/model"
	"go_zero/learn1/test2/user/rpc/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlxConn, c.CacheRedis),
	}
}
