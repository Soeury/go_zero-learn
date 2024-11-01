package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

// 这里是从配置文件中映射

// 默认只有 rest.RestConf，内部包含了 log，所以不需要手动嵌入
// 这里我们添加了 mysql 配置
// JWT 配置
// 和 mysql + cache 的配置
type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	Mysql struct {
		DataSource string
	}

	CacheRedis cache.CacheConf
}
