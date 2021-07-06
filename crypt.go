package hmapi

import (
	v1 "github.com/hlib-go/hcipher/v1"
	v2 "github.com/hlib-go/hcipher/v2"
	v3 "github.com/hlib-go/hcipher/v3"
)

// BizEncode 报文加密
func BizEncode(appid, biz string) (v string, err error) {
	if biz == "" {
		return
	}
	app, err := GetOpenApp(appid, true)
	if err != nil {
		return
	}
	switch app.Version {
	case APP_CEYPT_V3:
		v, err = v3.Encode(biz, app.MerPubKey) // 商户公钥加密
	case APP_CEYPT_V2:
		v, err = v2.EnAesRsa(v, app.MerPubKey) // 商户公钥加密
	case APP_CEYPT_V1:
		v, err = v1.EnDesMd5(biz, app.DesKey)
	default:
		v, err = v1.EnDesMd5(biz, app.DesKey)
	}
	if err != nil {
		return
	}
	return
}

var _rc = make(map[string]bool)

// BizDecode 报文解密
func BizDecode(appid, biz string) (v string, err error) {
	if biz == "" {
		return
	}
	app, err := GetOpenApp(appid, !_rc[appid])
	if err != nil {
		return
	}
	switch app.Version {
	case APP_CEYPT_V3:
		v, err = v3.Decode(biz, app.RsaPriKey) // 平台私钥解密
	case APP_CEYPT_V2:
		v, err = v2.DeAesRsa(v, app.RsaPriKey) // 平台私钥解密
	case APP_CEYPT_V1:
		v, err = v1.DeDesMd5(biz, app.DesKey)
	default:
		v, err = v1.DeDesMd5(biz, app.DesKey)
	}
	if err != nil && !_rc[appid] {
		_rc[appid] = true
		v, err = BizDecode(appid, biz)
		_rc[appid] = false
	}
	return
}

// BizEncodeT 报文加密
func BizEncodeT(appid, biz string) (v string, err error) {
	if biz == "" {
		return
	}
	app, err := GetOpenApp(appid, true)
	if err != nil {
		return
	}
	switch app.Version {
	case APP_CEYPT_V3:
		v, err = v3.Encode(biz, app.RsaPubKey) // 平台公钥加密
	case APP_CEYPT_V2:
		v, err = v2.EnAesRsa(v, app.RsaPubKey) // 平台公钥加密
	case APP_CEYPT_V1:
		v, err = v1.EnDesMd5(biz, app.DesKey)
	default:
		v, err = v1.EnDesMd5(biz, app.DesKey)
	}
	if err != nil {
		return
	}
	return
}
