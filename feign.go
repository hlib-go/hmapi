package hmapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hlib-go/hhttp"
	"io/ioutil"
	"net/http"
	"strings"
)

func Request(ctx context.Context, method, url string, body string) ([]byte, error) {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	var client = hhttp.Client()

	if ctx != nil {
		if ctx.Value("token") != nil {
			request.Header.Set("token", ctx.Value("token").(string))
		}
		appid := ctx.Value("appid")
		if appid != nil && appid != "" {
			request.Header.Set("appid", appid.(string))
		}
		ctxClient := ctx.Value("client")
		if ctxClient != nil {
			client = ctxClient.(*http.Client)
		}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusOK {
		return ioutil.ReadAll(response.Body)
	}

	// HTTP 返回非200状态都为异常
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if bytes != nil {
		var ferr *FeignErr
		err = json.Unmarshal(bytes, &ferr)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ferr.Errno + ":" + ferr.Error)
	}
	return nil, errors.New("99999:HTTP异常(" + response.Status + ")，请检查业务服务 " + method + " " + url)
}

type FeignErr struct {
	Errno string `json:"errno"`
	Error string `json:"error"`
}

func Do(ctx context.Context, method string, body string) ([]byte, error) {
	if method == "" {
		return nil, errors.New("接口名称不能为空")
	}
	name := strings.Split(method, ".")[1]
	if method[0:1] != "/" {
		method = "/" + method
	}
	url := "http://" + name + method
	return Request(ctx, "POST", url, body)
}

func Call(ctx context.Context, method string, i interface{}, o interface{}) (err error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return err
	}
	resp, err := Do(ctx, method, string(bytes))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, o)
	return err
}

// 内部接口请求
func InternalRequest(ctx context.Context, method string, path string, i interface{}, o interface{}) (err error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return err
	}
	resp, err := Request(ctx, method, "http://"+path, string(bytes))
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, o)
	return err
}
