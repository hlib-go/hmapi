package hmapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hlib-go/hmapi/errs"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ERR_HTTP_METHOD = errors.New("99999:请使用HTTP POST请求")
)

type ResolveParams struct {
	Context   context.Context
	Request   *http.Request
	Writer    http.ResponseWriter
	Body      []byte
	RequestId string
	Appid     string
}

func (p *ResolveParams) BodyUnmarshal(i interface{}) {
	err := json.Unmarshal(p.Body, i)
	if err != nil {
		panic(err)
	}
}

// 定义接口，注意：pattern必须/开头
func DefApi(pattern string, resolve func(p *ResolveParams) (out interface{}, err error)) {
	if pattern == "" || resolve == nil {
		panic(errors.New("DefApi()参数不能为空"))
		return
	}
	if !strings.HasPrefix(pattern, "/") {
		panic(errors.New("ERROR:***接口名必须已【/】开通"))
	}

	http.Handle(pattern, post(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			begTime   = time.Now().UnixNano() / 1e6
			header    = request.Header
			requestId = header.Get("requestId")
			appid     = header.Get("appid")
			hlog      = log.WithField("method", pattern[1:]).WithField("requestId", requestId).WithField("appid", appid)
			reqBytes  []byte
			resBytes  []byte
		)

		defer func() {
			if e := recover(); e != nil {
				var err = recoverError(e)
				resBytes = respFail(err)
			}
			hlog.WithField("ms", strconv.FormatInt(time.Now().UnixNano()/1e6-begTime, 10)).Info("响应报文 " + string(resBytes))
			responseWriter(writer, 200, "application/json;charset=utf-8", resBytes)
		}()

		reqBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		hlog.Info("请求报文 " + string(reqBytes))

		// 调用接口业务方法
		i, e := resolve(&ResolveParams{Appid: appid, Body: reqBytes, Request: request, Writer: writer, RequestId: requestId})
		if e != nil {
			panic(e)
		}
		resBytes = respSuccess(i)
	})))
}

func responseWriter(writer http.ResponseWriter, status int, contentType string, data []byte) {
	if status != 200 {
		writer.WriteHeader(status)
	}
	writer.Header().Set("Content-Type", contentType)
	writer.Write(data)
}

func post(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.ToUpper(request.Method) != "POST" {
			responseWriter(writer, 500, "application/json;charset=utf-8", respFail(ERR_HTTP_METHOD))
			return
		}
		handle.ServeHTTP(writer, request)
	})
}

func respSuccess(out interface{}) []byte {
	var suc = fmt.Sprintf(`{"errno":"%s","error":"%s"}`, errs.SUCCESS.Code, errs.SUCCESS.Msg)
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

func respFail(err error) []byte {
	var (
		errno = "ERROR"
		error = err.Error()
	)
	es := strings.Split(error, ":")
	if len(es) > 1 {
		errno = strings.Join(es[0:1], "")
		error = strings.Join(es[1:], ":")
	}
	return []byte(fmt.Sprintf(`{"errno":"%s","error":"%s"}`, errno, error))
}

func recoverError(e interface{}) (err error) {
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
	return err
}
