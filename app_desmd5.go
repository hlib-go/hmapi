package hmapi

import (
	desmd5 "github.com/hlib-go/hcipher/v1"
	"github.com/hlib-go/hmapi/errs"
)

// 解密
func DeDesMd5(appid, v string) (value string, err error) {
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
	value, err = desmd5.DeDesMd5(v, app.DesKey)
	if err != nil {
		app, err = GetOpenApp(appid, false) // 如果解密出错，重新读取密钥重试一次
		if err != nil {
			return
		}
		value, err = desmd5.DeDesMd5(v, app.DesKey)
	}
	return
}

// 加密
func EnDesMd5(appid, biz string) (value string, err error) {
	app, err := GetOpenApp(appid, true)
	if err != nil {
		return
	}
	value, err = desmd5.EnDesMd5(biz, app.DesKey)
	return
}
