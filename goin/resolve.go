package goin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hlib-go/hmapi/errs"
	log "github.com/sirupsen/logrus"
	"strings"
)

func handlerResolve(resolve func(p *ResolveParams) (out interface{}, err error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestId = ctx.GetString("requestId")
		var err error
		defer func() {
			if e := recover(); e != nil {
				switch v := e.(type) {
				case error:
					err = e.(error)
				case *log.Entry:
					err = errors.New(e.(*log.Entry).Message)
				case string:
					err = errors.New(e.(string))
				default:
					err = errors.New(fmt.Sprintf("ERROR: %v ", v))
				}
			}
			if err != nil {
				RenderJSON(ctx, respFail(err, requestId)) // Fail
			}
		}()
		out, err := resolve(&ResolveParams{
			Context:   ctx,
			RequestId: requestId,
		})
		if err != nil {
			return
		}
		RenderJSON(ctx, respSuccess(out, requestId)) // Success
	}
}

func respSuccess(out interface{}, requestId string) []byte {
	suc := fmt.Sprintf(`{"errno":"%s","error":"%s","requestId":"%s"}`, errs.SUCCESS.Code, errs.SUCCESS.Msg, requestId)
	if out == nil {
		return []byte(suc)
	}
	bytes, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}
	if len(bytes) > 2 {
		bytes = append([]byte(suc[0:len(suc)-1]+","), bytes[1:]...)
	} else {
		bytes = []byte(suc)
	}
	return bytes
}

func respFail(err error, requestId string) []byte {
	msg := err.Error()
	msgSplit := strings.Split(msg, ":")
	if msg[5:6] == ":" {
		return []byte(fmt.Sprintf(`{"errno":"%s","error":"%s","requestId":"%s"}`, msgSplit[0:1], strings.Join(msgSplit[1:], ":"), requestId))
	}
	return []byte(fmt.Sprintf(`{"errno":"%s","error":"%s","requestId":"%s"}`, "ERROR", msg, requestId))
}
