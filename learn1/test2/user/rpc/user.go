package main

import (
	"context"
	"flag"
	"fmt"

	// 注意包的导入
	"go_zero/learn1/test2/user/rpc/internal/config"
	"go_zero/learn1/test2/user/rpc/internal/server"
	"go_zero/learn1/test2/user/rpc/internal/svc"
	"go_zero/learn1/test2/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

// 启动之前记得开 redis 和 consul.....
// 如果检查了所有错误都不能运行，请检查一下请求URL的 query 参数类型是否有误。
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserInfoServer(grpcServer, server.NewUserInfoServer(ctx))

		// 映射 rpc 服务，只有在服务模式为 dev 或者 test 模式下才行(去配置文件中更改)
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 服务注册到 consul
	_ = consul.RegisterService(c.ListenOn, c.Consul)
	defer s.Stop()

	// 服务端拦截器注册在这里
	s.AddUnaryInterceptors(UserServerInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

// 服务端拦截器
func UserServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	// 处理请求前:
	fmt.Println("server intecepter IN...")

	// 取出请求 ctx 中的 metadata
	md, ok := metadata.FromIncomingContext(ctx)

	// 请求 ctx 中没有传入 metadata
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "miss metadata")
	}

	// 打印 metadata (请求的metadata不定义值也会默认自带一些参数在里面，所以打印不会为空)
	fmt.Printf("%+v\n", md)

	// 根据 metadata 中的数据进行校验操作
	if md["token"][0] != "zero-order-api" {
		return nil, status.Errorf(codes.Unauthenticated, "invaild token")
	}

	resp, err := handler(ctx, req) // 这里应该是实际处理 rpc 请求

	// 处理请求后:
	fmt.Println("server intecepter OUT...")
	return resp, err
}
