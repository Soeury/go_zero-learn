package logic

import (
	"context"
	"time"

	"go_zero/learn1/test2/user/api/internal/svc"
	"go_zero/learn1/test2/user/api/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {

	// 1. 处理用户请求 - 拿到用户名和密码(参数校验的步骤已经在 handler 层做完了，通过 req 可以直接拿到用户传入的信息)
	username := req.Username
	password := req.Password

	// 2. 检验用户信息与数据库中存入的是否一致
	//    -2.1 分别检验用户名和密码
	//    -2.2 通过用户名拿到用户，再检验密码 √
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, username)

	// 先检查用户是否存在错误，再检查数据库查询错误
	if err != nil {
		if err == sqlx.ErrNotFound {
			logx.Errorw(
				"login: query username not exist",
				logx.Field("err", err),
			)

			return &types.LoginResponse{
				Message: "User Not exist",
			}, nil
		}

		logx.Errorw(
			"login: find one by username failed",
			logx.Field("err", err),
		)

		return nil, ErrQueryFailed
	}

	// 加密算法在 signupHandler.go 里面
	saltPassword := EncryptPasswordWithSalt(password)
	if user.Password != saltPassword {
		logx.Errorw(
			"login: userpassword is not equal with password in db",
			logx.Field("err", err),
		)

		return &types.LoginResponse{
			Message: "login: password not equal",
		}, ErrPasswordNotEqual
	}

	// 生成 JWT
	secret := l.svcCtx.Config.Auth.AccessSecret
	iat := time.Now().Unix()
	expire := l.svcCtx.Config.Auth.AccessExpire
	token, err := l.GetJWTToekn(secret, iat, expire, user.UserId)
	if err != nil {
		logx.Errorw(
			"login create token failed",
			logx.Field("err", ErrWrongCreateToken),
		)

		return nil, ErrWrongCreateToken
	}

	// 3. 数据一致则登陆成功，不一致登陆失败
	return &types.LoginResponse{
		Message:      "login success",
		AccessToken:  token,
		AccessExpire: int(iat + expire),   // token 过期时间
		RefreshAfter: int(iat + expire/2), // token 刷新时间
	}, nil
}

// GetJWTtoekn 生成 jwt token
// 参数： 密钥 ， 签发时间 ， 令牌有效期 ， 用户 id
func (l *LoginLogic) GetJWTToekn(secretKey string, iat, seconds, user_id int64) (string, error) {

	// 创建一个空的负载映射
	claims := make(jwt.MapClaims)

	// 设置JWT的过期时间 = 签发时间 + 有效期
	claims["exp"] = iat + seconds

	// 签发时间
	claims["iat"] = iat

	// 自定义负载字段，记录用户 ID
	claims["userId"] = user_id

	// 创建一个新的JWT令牌，使用HS256签名算法
	token := jwt.New(jwt.SigningMethodHS256)

	// 将定义好的负载设置到 jwt 上
	token.Claims = claims

	// 使用密钥进行签名，返回 token 字符串
	return token.SignedString([]byte(secretKey))
}
