package goin

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGoin(t *testing.T) {
	SetSecret("12345678")
	POST("/test", func(p *ResolveParams) (out interface{}, err error) {
		token := p.VerifyToken()
		log.Info(token.Uid)
		out = gin.H{}
		return
	}, Cors)
	POST("/test2", func(p *ResolveParams) (out interface{}, err error) {
		token := p.VerifyToken()
		log.Info(token.Uid)
		out = gin.H{}
		return
	}, Cors)
}
