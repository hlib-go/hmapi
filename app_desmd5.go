package hmapi

import (
	"context"
	"github.com/hlib-go/hdesmd5"
	"github.com/hlib-go/hmapi/errs"
	"sync"
)

// 解密
func AppDecryptDesMd5(appid, v string) (value string, err error) {
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
	value, err = hdesmd5.DecryptDesMd5(v, app.DesKey)
	if err != nil {
		app, err = GetOpenApp(appid, false) // 如果解密出错，重新读取密钥重试一次
		if err != nil {
			return
		}
		value, err = hdesmd5.DecryptDesMd5(v, app.DesKey)
	}
	return
}

// 加密
func AppEncryptDesMd5(appid, biz string) (value string, err error) {
	app, err := GetOpenApp(appid, true)
	if err != nil {
		return
	}
	value, err = hdesmd5.EncryptDesMd5(biz, app.DesKey)
	return
}

var (
	_apps sync.Map
)

func GetOpenApp(appid string, readCache bool) (r *HmOpenAppCryptoResult, err error) {
	if readCache {
		value, ok := _apps.Load(appid)
		if ok {
			r = value.(*HmOpenAppCryptoResult)
			return
		}
	}
	r, err = HmOpenAppCrypto(nil, &HmOpenAppCryptoParams{Appid: appid})
	if err != nil {
		return
	}
	_apps.Store(appid, r)
	return
}

func HmOpenAppCrypto(ctx context.Context, params *HmOpenAppCryptoParams) (result *HmOpenAppCryptoResult, err error) {
	err = Call(ctx, "hm.open.app.crypto", params, &result)
	return
}

type HmOpenAppCryptoParams struct {
	RequestId string `json:"requestId"`
	Appid     string `json:"appid"`
}

type HmOpenAppCryptoResult struct {
	RequestId string `json:"requestId"`
	Appid     string `json:"appid"`
	Name      string `json:"name"`
	DesKey    string `json:"desKey"`
	AesKey    string `json:"aesKey"`
	RsaPubKey string `json:"rsaPubKey"`
	RsaPriKey string `json:"rsaPriKey"`
}
