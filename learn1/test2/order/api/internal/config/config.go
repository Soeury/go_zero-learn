package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Mysql struct {
		DataSource string
	}

	UserRPC zrpc.RpcClientConf // 连接其他微服务的 rpc 客户端

	CacheRedis cache.CacheConf
}
