package hmapi

// 全局配置
var _options = new(Options)

type Options struct {
	TokenSecret string `json:"tokenSecret"` // token加密密钥
}

func SetOptions(o *Options) {
	_options = o
}

func SetTokenSecret(secret string) {
	_options.TokenSecret = secret
}
