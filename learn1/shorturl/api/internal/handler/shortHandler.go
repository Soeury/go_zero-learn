package handler

import (
	"net/http"

	"go_zero/learn1/shorturl/api/internal/logic"
	"go_zero/learn1/shorturl/api/internal/svc"
	"go_zero/learn1/shorturl/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShortHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShortLogic(r.Context(), svcCtx)
		resp, err := l.Short(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {

			// 这是原来返回的响应数据
			// httpx.OkJsonCtx(r.Context(), w, resp)

			// 现在我们需要进行重定向
			// (短链接实现原理)
			http.Redirect(w, r, resp.Long, http.StatusFound)
		}
	}
}
