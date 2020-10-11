package hapi

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
