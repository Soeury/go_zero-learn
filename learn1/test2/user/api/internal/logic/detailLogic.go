package logic

import (
	"context"
	"fmt"
	"strconv"

	"go_zero/learn1/test2/user/api/internal/svc"
	"go_zero/learn1/test2/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailRequest) (resp *types.DetailResponse, err error) {

	// jwt 鉴权后，从 JWT 中取出数据( l.ctx.Value())
	fmt.Printf("JWT get user_id:%v\n", l.ctx.Value("userId"))

	// 1. 用户参数校验
	//    ---- 这一步再上一层已经做好了
	user_id, _ := strconv.Atoi(req.UserID)

	// 2. 数据库查询数据
	user, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, int64(user_id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			logx.Errorw(
				"detail: query with no rows return",
				logx.Field("err", err),
			)

			return nil, err
		}

		logx.Errorw(
			"detail: find one by user_id failed",
			logx.Field("err", err),
		)

		return nil, ErrQueryFailed
	}

	// 3. 数据格式化(数据库中的数据和前端需求的数据不一致)
	// 4. 返回响应
	return &types.DetailResponse{
		Username: user.Username,
		Gender:   int(user.Gender),
	}, nil
}
