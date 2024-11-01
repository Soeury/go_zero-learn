package logic

import (
	"context"
	"strconv"

	"go_zero/learn1/test2/order/api/internal/errorx"
	"go_zero/learn1/test2/order/api/internal/interceptor"
	"go_zero/learn1/test2/order/api/internal/svc"
	"go_zero/learn1/test2/order/api/internal/types"
	"go_zero/learn1/test2/user/rpc/userinfo"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 日志的 err 就是用 err 本身 ， 需要 return 的 err 才需要进行包装，返回指定的 err
func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchResp, err error) {

	// 1. 根据 order_id 查询订单信息
	orderid := uint64(req.OrderID)
	order_info, err := l.svcCtx.UserModel.FindOneByOrderId(l.ctx, orderid)

	if err != nil {
		if err == sqlx.ErrNotFound {
			logx.Errorw(
				"order search: find info by order_id return no rows",
				logx.Field("err", err),
			)

			// 我们自定义了一个错误包，就不需要额外的错误定义了
			// return nil, ErrReqDataInvalid
			return nil, errorx.NewRespErr(errorx.CodeInvalidParam)
		}

		logx.Errorw(
			"order search: query db faield",
			logx.Field("err", err),
		)

		return nil, errorx.NewRespErr(errorx.CodeQueryDBFailed)
	}

	// 2. 根据拿到的订单信息中的 user_id 通过 rpc 调用查询用户信息
	// 先假设: 查询出来的  user_id  =  1729407198 , 之后写好了也用这个
	user_id := int64(order_info.UserId)

	// ctx 传入值，传入的值默认都是 any 类型
	// 如果存在名称相同的其他变量，会出现冲突
	// 解决办法: 自定义变量类型，键通过使用自定义类型声明的常量的方式传入 WithValue 函数
	l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdminId, "239514")

	userResp, err := l.svcCtx.UserRPC.GetUserInfo(l.ctx, &userinfo.UserInfoReq{UserId: user_id})
	if err != nil {
		logx.Errorw(
			"order search: get userinfo failed",
			logx.Field("err", err),
		)

		return nil, errorx.NewRespErr(errorx.CodeRPCFailed)
	}

	// 3. 按照要求拼接返回的数据
	return &types.SearchResp{
		UserName: userResp.GetUsername(),
		Status:   100,
		UserID:   strconv.Itoa(int(user_id)),
	}, nil
}
