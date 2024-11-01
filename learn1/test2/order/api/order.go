package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"go_zero/learn1/test2/order/api/internal/config"
	"go_zero/learn1/test2/order/api/internal/errorx"
	"go_zero/learn1/test2/order/api/internal/handler"
	"go_zero/learn1/test2/order/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul" // 这里需要匿名导入这个包
)

// "learn1/test2/order/api/etc/order-api.yaml"
// "C:/Users/Chenk/Desktop/go_zero/learn1/test2/order/api/etc/order-api.yaml"
// go run "c:\Users\Chenk\Desktop\go_zero\learn1\test2\order\api\order.go"
var configFile = flag.String("f", "etc/order-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 打印信息确认加载是否成功
	fmt.Printf("%+v\n", c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// go_zero 处理返回的错误
	// 类似于一个钩子函数
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		switch e := err.(type) {
		case *errorx.ResponseErr: // 返回的是我们定义好的错误
			return http.StatusOK, e.Data()
		default: // 返回的不是我们定义好的错误
			return http.StatusInternalServerError, nil
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
