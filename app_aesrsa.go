package hmapi

import (
	aesrsa "github.com/hlib-go/hcipher/v2"
	"github.com/hlib-go/hmapi/errs"
)

// 加密
func EnAesRsa(appid, v string) (value string, err error) {
	app, err := GetOpenApp(appid, true)
	if err != nil {
		return
	}
	value, err = aesrsa.EnAesRsa(v, app.RsaPubKey)
	return
}

// 解密
func DeAesRsa(appid, v string) (value string, err error) {
	if appid == "" {
		err = errs.E99901
		return
	}
	if v == "" {
		return
	}
	app, err := GetOpenApp(appid, true)
	if err != nil {
		return
	}
	// 解密
	value, err = aesrsa.DeAesRsa(v, app.RsaPriKey)
	if err != nil {
		app, err = GetOpenApp(appid, false) // 如果解密出错，重新读取密钥重试一次
		if err != nil {
			return
		}
		value, err = aesrsa.DeAesRsa(v, app.RsaPriKey)
	}
	return
}
