package errorx

// 自定义的响应码 and 错误包

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeRPCFailed
	CodeDefault
	CodeInvalidParam
	CodeInvalidToken
	CodeInvalidPassword
	CodeQueryDBFailed
	CodeTargetExist
	CodeUserNotExist
	CodeServerBusy
	CodeNeedLogin
	CodeInsertDBFailed
)

// 这里比老师写的好一点的地方就是不需要手动传入字符串了......
var CodeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "wrong query params",
	CodeTargetExist:     "target exist",
	CodeUserNotExist:    "user not exist",
	CodeInvalidPassword: "wrong password",
	CodeServerBusy:      "busy server",
	CodeNeedLogin:       "need login",
	CodeInvalidToken:    "wrong token",
	CodeQueryDBFailed:   "query db failed",
	CodeRPCFailed:       "get RPC server faield",
	CodeInsertDBFailed:  "insert into db faield",
}

// Msg 状态码调用 msg 方法可以得到状态码对应的错误信息
func (p *ResponseErr) Data() *ResponseData {

	return &ResponseData{
		Code: p.Code,
		Msg:  p.Msg,
	}
}

// 错误 结构体 自定义
type ResponseErr struct {
	Code ResCode `json:"code"`
	Msg  string  `josn:"msg"`
}

// 响应 结构体 自定义
type ResponseData struct {
	Code ResCode `json:"code"`
	Msg  string  `josn:"msg"`
}

// 结构体重写 error 接口
func (re *ResponseErr) Error() string {

	return re.Msg
}

// 返回的是错误类型，不需要传入额外的字符串，只需要传入对应的状态码即可
func NewRespErr(code ResCode) error {

	// 这里返回的是指针类型还是非指针类型决定了后面处理的错误的断言的类型
	return &ResponseErr{
		Code: code,
		Msg:  CodeMsgMap[code],
	}
}

// ? ? ? 这个的用途何在 ? ? ?
// 返回默认的错误 - 系统繁忙
func NewRespDefaultErr() error {

	return &ResponseErr{
		Code: CodeServerBusy,
		Msg:  CodeMsgMap[CodeServerBusy],
	}
}

// ----------------------------------------------------------------------------------------

/*

// 自定义的错误
type CodeErr struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 自定义的错误响应
type CodeErrResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// NewCodeErr 返回自定义错误
func NewCodeErr(code int, msg string) error {

	return &CodeErr{
		Code: code,
		Msg:  msg,
	}
}

// NewDefaultCodeErr 返回默认的自定义错误
func NewDefaultCodeErr(msg string) error {

	return &CodeErr{
		// Code: DefauleCodeErr,
		Msg: msg,
	}
}

// Error  CodeErr 结构体重写 error 接口
func (e CodeErr) Error() string {

	return e.Msg
}

// Data  返回自定义类型的错误响应
func (e *CodeErrResponse) Data() *CodeErrResponse {

	return &CodeErrResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

*/
