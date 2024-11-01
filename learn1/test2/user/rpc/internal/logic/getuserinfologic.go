package logic

import (
	"context"

	"go_zero/learn1/test2/user/rpc/internal/svc"
	"go_zero/learn1/test2/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.UserInfoReq) (*user.UserInfoResp, error) {

	// 注意这个变量名字不要和包名字重复了
	one, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.GetUserId())
	if err != nil {
		if err == sqlx.ErrNotFound {
			logx.Errorw(
				"rpc getUserInfo by id returns no rows",
				logx.Field("err", err),
			)

			return nil, err
		}

		logx.Errorw(
			"rpc getUserInfo by id failed",
			logx.Field("err", err),
		)

		return nil, ErrQueryFailed
	}

	return &user.UserInfoResp{
		UserId:   one.UserId,
		Username: one.Username,
		Gender:   one.Gender,
	}, nil
}
