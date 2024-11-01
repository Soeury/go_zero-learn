package main

import (
	"flag"
	"fmt"

	"go_zero/learn1/test2/user/api/internal/config"
	"go_zero/learn1/test2/user/api/internal/handler"
	"go_zero/learn1/test2/user/api/internal/middleware"
	"go_zero/learn1/test2/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	// 1. 从配置文件加载信息
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 2. 打印信息确认是否加载成功
	fmt.Printf("%+v\n", c)

	// 3. ? ? ? 好像是加载默认配置
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// x. 使用全局中间件
	server.Use(middleware.CopyResponse)

	// 4. 返回配置文件结构体(这个很重要)，并传入到 handler 包中进行使用
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
