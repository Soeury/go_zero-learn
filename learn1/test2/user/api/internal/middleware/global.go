package middleware

import (
	"bytes"
	"fmt"
	"net/http"
)

// 定义结构体，实现了 http.ResponseWriter 接口
type BodyCopy struct {
	http.ResponseWriter // 结构体嵌入接口，就试实现了接口中的所有方法，通过重写的方式可以给这些方法的功能
	Body                *bytes.Buffer
}

// NewBodyCopy 初始化结构体
func NewBodyCopy(w http.ResponseWriter) *BodyCopy {

	return &BodyCopy{
		ResponseWriter: w,
		Body:           bytes.NewBuffer([]byte{}),
	}
}

// Write 结构体重写接口的 write 方法
func (bc *BodyCopy) Write(b []byte) (int, error) {

	// 将响应数据记录到缓存中
	// 这个 write 是 buffer 自己的 write
	bc.Body.Write(b)

	// 然后往 http 响应中写入响应内容
	// 这个 write 是 http.ResponseWriter 中的 write，作用是: 往响应中写入数据
	return bc.ResponseWriter.Write(b)
}

// CopyResponse 自定义中间件
func CopyResponse(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// 初始化得到一个自定义的 http.ResponseWriter
		bc := NewBodyCopy(w)
		next(w, r)
		// 处理请求后，可以拿到请求数据和响应数据啦
		fmt.Printf("req:%v , resp:%v\n", r.URL, bc.Body.String())
	}
}
