package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type CostMiddleware struct {
}

func NewCostMiddleware() *CostMiddleware {
	return &CostMiddleware{}
}

func (m *CostMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		now := time.Now()
		next(w, r)                                      // 执行后续的 handler 函数
		fmt.Printf("--->cost:%d", int(time.Since(now))) // 打印执行时间的中间件
	}
}
