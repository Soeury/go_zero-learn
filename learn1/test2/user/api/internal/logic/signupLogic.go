package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"go_zero/learn1/test2/user/api/internal/svc"
	"go_zero/learn1/test2/user/api/internal/types"
	"go_zero/learn1/test2/user/model"

	"github.com/zeromicro/go-zero/core/logx" // go_zero 自带的 logx 包内部封装了 logger 接口
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const Salt = "boliang and rabbit"

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// EncryptPasswordWithSalt 使用了 salt 来“盐化”密码
func EncryptPasswordWithSalt(opassword string) string {

	combined := fmt.Sprintf("%s%s", opassword, Salt)
	h := md5.New()
	h.Write([]byte(combined))

	// 传入 nil 表示我们不需要提前分配一个切片
	return hex.EncodeToString(h.Sum(nil))
}

// NewSignupLogic 这个 NewStruct 是供处理器 Handler 调用的，处理器拿到结构体对象后，就可以调用对应的结构体方法
func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {

	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Signup 这里写我们的业务逻辑
func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SignupResponse, err error) {

	// 参数校验(其实不应该写在这里,后面会转移)
	if req.Password != req.RePassword {
		return nil, ErrPasswordNotEqual
	}

	logx.Debugv(req)
	logx.Debugf("req:%+v\n", req)

	// 通过调用对象 SignupLogic 中的  svcCtx 可以拿到 UserModel 中的方法，然后按需调用即可

	//  -0. 查询 username 是否已经存在? 有三种返回的结果,可以进行错误处理
	//      -0.1. 查询数据库失败
	//      -0.2. 没有返回的记录，这里就继续往下执行
	//      -0.3. 有返回记录
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && err != sqlx.ErrNotFound {

		// 日志记录
		logx.Errorw(
			"sign up find one by user name failed",
			logx.Field("err", err),
		)
		fmt.Printf("find one by username err:%s\n", err)
		return nil, ErrQueryFailed
	}

	if u != nil {
		return nil, ErrUserExist
	}

	//  -1. 生成 userid (雪花算法,这里简略使用时间戳来作为ID,后面会重新写)
	//  -2. 加密 password (加盐 md5...)
	newPassword := EncryptPasswordWithSalt(req.Password)

	userData := &model.User{
		UserId:   time.Now().Unix(),
		Username: req.Username,
		Password: newPassword,
		Gender:   int64(req.Gender),
	}

	//  -3. 将用户信息保存到数据库中
	_, err = l.svcCtx.UserModel.Insert(context.Background(), userData)
	if err != nil {
		logx.Errorf("signup insert failed:%v\n", err)
		return nil, err
	}

	return &types.SignupResponse{
		Message: "sign up success!",
	}, nil
}
