package goin

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/hlib-go/hgenid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"time"
)

const json_resp_body = "middleware_json_resp_body"

// 使用此方法返代替ctx.JSON() 才能记录到日志
func RenderJSON(ctx *gin.Context, data []byte) {
	ctx.Header("Content-Type", "application/json;charset=UTF-8")
	ctx.Set(json_resp_body, data)
	ctx.Data(200, "application/json;charset=UTF-8", data)
}

// 中间件：记录接口请求响应报文日志
func Logger(ctx *gin.Context) {
	begTime := time.Now().UnixNano()

	// 读取请求ID
	requestId := ctx.Query("requestId")
	if requestId == "" {
		requestId = hgenid.UUID()
	}
	ctx.Set("requestId", requestId)

	rlog := log.WithField("requestId", requestId).WithField("path", ctx.FullPath())

	// 请求
	{
		rlog.Info("请求方式 ", ctx.Request.Method, " ", ctx.Request.RequestURI)
		if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" {
			bodyBytes, err := ctx.GetRawData()
			if err != nil {
				rlog.Error("ctx.GetRawData Error " + err.Error())
				ctx.Abort()
				return
			}
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 取值后重新复制，否则后续无法取值
			rlog.Info("请求报文 ", string(bodyBytes))
		}
	}

	ctx.Next()

	// 响应
	{
		respBody, ok := ctx.Get(json_resp_body)
		if !ok {
			respBody = "{}"
		}
		rlog.WithField("ms", strconv.FormatInt((time.Now().UnixNano()-begTime)/1e6, 10)).Info("响应报文 ", respBody)
	}

}
