package goin

import (
	"github.com/gin-gonic/gin"
	"github.com/hlib-go/hmapi/token"
)

type Goin struct {
	Engine *gin.Engine
	Secret string // token加密密钥 ,在使用token前通过 goin.SetSecret(v) 设置
}

var _g = &Goin{
	Engine: gin.Default(),
}

func GetEngine() *gin.Engine {
	return _g.Engine
}

func SetEngine(e *gin.Engine) {
	_g.Engine = e
}

func SetSecret(secret string) {
	_g.Secret = secret
}

func Use(handlers ...gin.HandlerFunc) {
	_g.Engine.Use(handlers...)
}

type ResolveParams struct {
	Context   *gin.Context
	RequestId string
	Token     *token.Object
}

func (p *ResolveParams) BindJSON(obj interface{}) {
	if err := p.Context.BindJSON(obj); err != nil {
		panic(err)
	}
}

func (p *ResolveParams) VerifyToken() *token.Object {
	t, err := token.Ver(_g.Secret, p.Context.GetHeader("token"))
	if err != nil {
		panic(err)
	}
	p.Token = t
	return p.Token
}

// POST 请求响应都为JSON数据 ， 中间件统一日志记录方式、错误处理方式
func POST(relativePath string, resolve func(p *ResolveParams) (out interface{}, err error), handlers ...gin.HandlerFunc) {
	if handlers == nil {
		handlers = []gin.HandlerFunc{}
	}
	handlers = append(
		handlers,
		Logger,
		handlerResolve(resolve))
	_g.Engine.POST(relativePath, handlers...)
}
