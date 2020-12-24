package hmapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 全局错误码定义
// 999**
var (
	SUCCESS = New("00000", "ok")
	FAIL    = New("99999", "fail")

	E99901 = New("99901", "无效appid")
	E99902 = New("99902", "解密失败")
	E99903 = New("99902", "加密失败")

	E99910 = New("99910", "无效code")
	E99911 = New("99911", "无效token")
)

type Err struct {
	Code string `json:"errno"`
	Msg  string `json:"error"`
}

func New(code, msg string) *Err {
	return &Err{Code: code, Msg: msg}
}

func NewF(err error) *Err {
	e := err.Error()
	i := strings.Index(e, ":")
	if i == -1 {
		return &Err{Code: "ERROR", Msg: e}
	}
	return &Err{Code: e[0:i], Msg: e[i+1:]}
}

func (e *Err) NewMsg(msg string) *Err {
	return New(e.Code, msg)
}

func (e *Err) NewMsgF(args ...interface{}) *Err {
	return New(e.Code, fmt.Sprintf(e.Msg, args...))
}

func (e *Err) Error() string {
	return e.Code + ":" + e.Msg
}

func (e *Err) JsonMarshal() []byte {
	b, _ := json.Marshal(e)
	return b
}
