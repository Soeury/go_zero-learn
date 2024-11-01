## info
info(
    title: "user"
    desc: "how to write api file"
    author: "rabbit"
    email: @kedudu.com
    version: 1.0
) 





## .api 文件自动生成结构文件
1. 编写完 .api 文件后，切换路径，然后终端输入命令:  
    goctl api go -api name.api -dir . -style=goZero

2. 然后会在当前目录下生成 etc 和 internal 文件夹，里面会放置相关内容
    - 2.1. etc 文件夹内置 .yaml 配置文件的内容
    - 2.2. internal 文件夹内置 读取配置文件，请求处理器，请求路由，业务逻辑......

3. internal 包一般是不可以被外部导入的包(规定) 





## JWT 进行用户鉴权
1. 用户登陆后，生成 jwt，放在响应数据中返回

2. 后续用户的每一次请求都会带上 token , 后端从请求头中获取 token 进行解析
    - 2.1 解析成功后，用户则登录
    - 2.2 解析失败或者token 过期，用户不登陆 





## @server(
        handler:   指定路由，进行以下操作
            jwt:   开启jwt认证
     middleware:   添加中间件
          group:   指定的路由都会分到这个组里
        timeout:   超时配置
         perfix:   路由前缀
       maxBytes:   添加请求体大小控制
)





## API自定义中间件
1. .api 文件中在 @server() 内部添加 middleware 选项，并定义中间件的名字

2. 在生成的 middleware 文件中实现中间件的功能

3. 在 srv 文件中的 ServiceContext 和 NewServiceContext 添加中间件的选项，注意字段名要和.api文件中保持一致





## 全局中间件： 举例：实现记录所有的响应数据的功能
1. 要求： 只要在中间件那一层按照格式定义中间件函数即可

2. 应用： 在 main 函数中使用 server.Use(middleware_name) 全局注册即可  

    - 2.1  type middleware func(next http.HandlerFunc) http.HandlerFunc
    - 2.2  type Handle func(ResponseWriter, *Request)

3. 中间件框架
   func CopyResponse(next http.HandlerFunc) http.HandlerFunc {
  	 return func(w http.ResponseWriter, r *http.Request) {
 	 	  // 处理请求前
 		  next(w, r)
 		  // 处理请求后
 	     }
   }





## mysql with no cache
---------- 创建表 ----------

1. 在 name.sql 中放置创建表的语句

2. 在数据库中创建表

3. 终端输入命令自动生成 crud 的代码块:
    goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/zero" -table="user" -dir=./model 

---------- 更改配置文件 ----------

4. api/config/config.go 增加一个字段:  Mysql struct { 
                                            DataSource string 
                                        }

5. api/etc/user-api.yaml 增加 mysql 对应的内容: Mysql:
                                                root:123456@tcp(127.0.0.1:3306)/zero?parseTime=True&loc=Local

6. api/internal/svc 文件的更改：  
        -1. ServiceContext  结构体  NameModel model.interface
        -2. NewServiceContext 函数  NameModel : model.NewInterface(sqlxConn)
        -3. sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)  // 这个连接是 go_zero 封装好了的

---------- 业务逻辑编写 ----------

7. 业务逻辑调用结构体中新增接口的对应的方法





## mysql with cache
1. 终端输入： goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/zero" -table="user" -dir=./model -c

2. 其他的同 with no cache ，需要修改的如下:
	//  ------- with cache -------
	//  创建的如果是携带缓存版本的代码
	//  首先需要修改 etc 文件和 config 文件
	//  etc 文件注意，一个 - Host 代表以台主机，这台主机的其余参数(password 等，如果有的话)不需要携带 -
	//  下面的 UserModel 的函数需要两个参数





## go_zero RPC
1.  编写 protobuffer 文件

2.  切换路径到 rpc 文件下面! ! 

3.  终端输入命令:  goctl rpc protoc user.proto --go_out=./types --go-grpc_out=./types --zrpc_out=./

4.  修改 config 文件和 yaml 文件
    yaml 文件中的ETCD 配置:
        Etcd:
        Hosts:
        - 127.0.0.1:2379
        Key: user.rpc

5.  修改 srv 文件中的信息(与config对应)

6.  编写业务代码

7. 检验代码，这里为了方便，使用了第三方库的 grpc ui 界面来检查服务
    - 7.1  地址：go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
    - 7.2  步骤： 启动服务后，终端输入: grpcui -plaintext localhost:port , 之后根据 url 在浏览器中打开即可
    - 7.3  可能会出现这个错误： 
            failed to compute set of methods to expose : server does not support the reflection API 
         更改 rpc 服务运行模式即可





## order 通过 RPC 调用 user 服务
1. /api/order/search : 根据 id 查询订单信息

2. 更改 order config , etc 配置
    - 2.1 config: UserRpc zrpc.RpcClientConf
    - 2.2 etc: 自己看着办

3. 修改 srv 结构体字段
    - 3.1 go-zero 会自动生成 GRPC 客户端
    - 3.2 修改 svctx 结构体: 
        - UserRPC userinfo.UserInfo  这个类型表示 需要调用的 rpc 服务的包下面的 inteface 接口
    - 3.3 修改 NewSvctx 函数:
        - UserRPC: userinfo.NewUserInfo(zrpc.MustNewClient(c.UserRPC))  go_zero 将客户端连接已经生成好，只需要传入配置





## 使用 conusl 作为注册中心

## ------------ 服务注册 ------------
1. 修改 config 文件
    - 1.1 导入包：go get -u github.com/zeromicro/zero-contrib/zrpc/registry/consul
    - 1.2 config 配置中添加： Consul consul.Conf

2. 修改 etc 文件 Consul 添加  host , key 选项

3. 服务注册到 consul：  主函数添加： consul.RegisterService(c.ListenOn, c.Consul)

## ------------ 服务发现（调用服务order的一方去修改配置文件） ------------
1. 只修改 etc 文件， config 中的配置只在 zrpc.RpcClientConf  中
    - 1.1 etc 文件给  UserRPC 添加 Traget: consul://127.0.0.1:8500/zero_user_rpc?wait=30s 字段
    - 1.2 调用服务的一方主函数匿名导入包：_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
    (这里老师讲过，服务发现的两种方式，一是手动写代码进行服务注册，二是在创建客户端时的地址改成上面的,要求是要匿名导入包)





## 这里安装 zero_contrib 和 go_zero 的时候 可能会出一些问题
// 都是包版本的问题，哪个包出问题了就把他更新下载到最新版本即可
// 比如现在遇到的 go.opentelemetry.io/otel/sdk@v1.24.0 就出问题了
// 解决： 去 go.mod 里面， ctrl 点击包，进入浏览器，找到最新版本下载之后 vscode 会自动替换，之后就没问题了
// 遇到问题，看终端出现了什么问题，找解决方法就行了





## GRPC 调用服务传递 metadata
1. 什么是 metadata ? 有什么用 ? 
2. grpc 客户端和服务端的拦截器 ? 




## go_zero  项目添加 rpc client 拦截器
- order.search 服务作为 rpc 客户端 
1. client拦截器 函数：在 创建 rpc 客户端时加进去( Newsrv 函数里面:  zrpc.WithUnaryClientInterceptor(拦截器名字))
    func Name(
        ctx context.Context,
        method string,
        req, reply interface{},
        cc *grpc.ClientConn,
        invoker grpc.UnaryInvoker,
        opts ...grpc.CallOption,
    ) error {}

2. 这个拦截器函数怎么添加外部数据呢？通过 ctx 取值然后传入 metadata ，metadata 编写完成之后，传入 ctx 中
    - ctx = metadata.NewOutgoingContext(ctx, md)

3. 添加 metadata 数据使用包:  "google.golang.org/grpc/metadata"

4. ctx 什么时候存入值呢？在发起 rpc 调用之前通过 context.WithValue 传入值(logic 层)
	- ctx 传入值，传入的值默认都是 any 类型
	- 如果存在名称相同的其他变量，会出现冲突
	- 解决办法: 自定义变量类型，键通过使用自定义类型声明的常量的方式传入 WithValue 函数

    - type CtxKey string
    - const CtxKeyAdminId CtxKey = "adminID"




## go_zero  项目添加 rpc server 拦截器
- user rpc 作为 rpc server 
1. server 拦截器添加在哪里？ 添加在 rpc server 的 main 函数中:  s.AddUnaryInterceptors(拦截器名字)

2. server 拦截器:
    func UserServerInterceptor(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {}

3. ! ! ! ! ! ! 为什么不加拦截器可以正常运行，加了拦截器不能正常运行 ?





## 错误处理
- 日志的 err 就是用 err 本身 ， 需要 return 的 err 才需要进行包装，返回指定的 err

1. 自定义错误格式： 见 order 中的 errorx 包

2. 业务逻辑中使用定义好的错误(日志记录的 err 用 err 本身 ， 需要 return 的 err 才需要进行包装，返回我们指定的 err)

3. go_zero 处理返回的错误:
    - 3.1 errorx 包在哪个服务的哪个模块下，处理的函数就写在哪个模块的 main 函数中





## go_zero 框架中的 goctl 模板文件
1. 模板的用处: goctl 指令根据模板来生成代码

```bash
    goctx api go -api name.api -dir . -style=goZero
```

2. goctl 模板:

    - goctl template -h      查看 goctl 支持的命令
    - goctl env              查看模板文件的位置... ( 第一次查看的时候可能不存在模板文件)
    - goctl template init    初始化模板文件 

3. 初始化模板文件:
```bash
    goctl tempalte init
```

4. 要点总结: 
    - 1. 熟悉 go 的模板语法
    - 2. 知道 goctl 的模板文件存放在哪里
    - 3. 具体使用:
        - 1. 找出模板文件并按需修改
        - 2. 生成代码 (存在同名文件就不会生成)
⭐ 去学一下怎么写⭐
