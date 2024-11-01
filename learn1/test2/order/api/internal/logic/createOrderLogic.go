package logic

import (
	"context"
	"strconv"
	"time"

	"go_zero/learn1/test2/order/api/internal/errorx"
	"go_zero/learn1/test2/order/api/internal/svc"
	"go_zero/learn1/test2/order/api/internal/types"
	"go_zero/learn1/test2/order/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// order_id ~= username 不能重复,需要用户传入    trade_id ~= userID  (自动生成不重复的)，不需要传入
func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (resp *types.CreateOrderResp, err error) {

	// 1. 参数校验
	if req.UserID == 0 {
		logx.Infow(
			"order create with nil id",
			logx.Field("err", errorx.NewRespErr(errorx.CodeInvalidParam)),
		)

		return nil, errorx.NewRespErr(errorx.CodeInvalidParam)
	}

	// 打印(非必须)
	logx.Debugv(req)
	logx.Debugf("req:%+v\n", req)

	// 2. 查询订单 id 是否已经存在
	//    -2.1. 数据库错误
	//    -2.2. 已经存在数据
	//    -2.3. 不存在数据

	order, err := l.svcCtx.UserModel.FindOneByOrderId(l.ctx, uint64(req.OrderID))
	if err != nil && err != sqlx.ErrNotFound {

		// 日志记录
		logx.Errorw(
			"order create find one by order id failed",
			logx.Field("err", err),
		)

		return nil, errorx.NewRespErr(errorx.CodeQueryDBFailed)
	}

	// 订单存在
	if order != nil {
		return nil, errorx.NewRespErr(errorx.CodeTargetExist)
	}

	user_id := uint64(req.UserID)
	order_id := uint64(req.OrderID)
	trade_id := strconv.Itoa(int(time.Now().Unix() + 1234567890))
	data := &model.Order{
		UserId:  user_id,
		OrderId: order_id,
		TradeId: trade_id,
	}

	// 订单不存在，数据保存到数据库中
	_, err = l.svcCtx.UserModel.Insert(l.ctx, data)
	if err != nil {
		logx.Errorw(
			"order create insert field",
			logx.Field("err", err),
		)
		return nil, errorx.NewRespErr(errorx.CodeInsertDBFailed)
	}

	// 返回响应
	return &types.CreateOrderResp{
		Message: "Success!",
	}, nil
}
