package hmapi

import (
	"context"
	"github.com/hlib-go/hmapi/errs"
	"sync"
)

var (
	_apps sync.Map
)

func GetOpenApp(appid string, readCache bool) (r *HmOpenAppCryptoResult, err error) {
	if appid == "" {
		err = errs.E99901
		return
	}
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
	RequestId string     `json:"requestId"`
	Appid     string     `json:"appid"`
	Name      string     `json:"name"`
	DesKey    string     `json:"desKey"`
	AesKey    string     `json:"aesKey"`
	RsaPubKey string     `json:"rsaPubKey"`
	RsaPriKey string     `json:"rsaPriKey"`
	Version   AppVersion `json:"version"`
	MerPubKey string     `json:"merPubKey"`
	MerPriKey string     `json:"merPriKey"`
}

type AppVersion string

const (
	APP_CEYPT_V1 AppVersion = "v1" // des +md5
	APP_CEYPT_V2 AppVersion = "v2" // aes + rsa
	APP_CEYPT_V3 AppVersion = "v3" // des+sha1+rsa
)
