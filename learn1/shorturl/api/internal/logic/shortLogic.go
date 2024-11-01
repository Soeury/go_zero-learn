package logic

import (
	"context"

	"go_zero/learn1/shorturl/api/internal/svc"
	"go_zero/learn1/shorturl/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLogic {
	return &ShortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShortLogic) Short(req *types.Request) (resp *types.Response, err error) {

	if req.Short == "rabbit" {
		return &types.Response{Long: "http://www.liwenzhou.com"}, nil
	}

	return &types.Response{Long: "http://www.baidu.com"}, nil
}
