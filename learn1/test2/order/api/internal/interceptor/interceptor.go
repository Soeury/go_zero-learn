package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 给 ctx 的键自定义一种类型，避免名称冲突
type CtxKey string

const (
	CtxKeyAdminId CtxKey = "adminID"
)

// 定义 order rpc 客户端拦截器
func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {

	// rpc 调用前: 编写 客户端的调用逻辑
	fmt.Println("client intecepter IN...")

	// 这个拦截器函数怎么添加外部数据 ?  ctx 取值然后传入。 ctx 什么时候存入值呢？在发起 rpc 调用之前(logic 层)
	adminID := ctx.Value(CtxKeyAdminId).(string)
	md := metadata.Pairs(
		"key1", "value1", // "key1" 的值将会是:  []string{"value1" , "value1"}
		"key1", "value2",
		"token", "zero-order-api",
		"adminID", adminID,
	)

	// 将元数据加入到 context 中
	ctx = metadata.NewOutgoingContext(ctx, md)
	err := invoker(ctx, method, req, reply, cc, opts...) // 这一步相当于 rpc 调用

	// rpc 调用后:
	fmt.Println("client intecepter OUT...")
	return err
}
