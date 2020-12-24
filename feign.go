package hmapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hlib-go/hgenid"
	"github.com/hlib-go/hhttp"
	"io/ioutil"
	"net/http"
	"strings"
)

// 内网调用接口
/*type Feign interface {
	Do(ctx context.Context, method string, body interface{}, result interface{}) (err error)
}*/

type Error500 struct {
	Errno string `json:"errno"`
	Error string `json:"error"`
}

func Post(ctx context.Context, url string, body string) ([]byte, error) {
	request, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	var client = hhttp.Client()

	if ctx != nil {
		appid := ctx.Value("appid")
		if appid != nil && appid != "" {
			request.Header.Set("appid", appid.(string))
		}
		ctxClient := ctx.Value("client")
		if ctxClient != nil {
			client = ctxClient.(*http.Client)
		}
	}
	request.Header.Set("tid", hgenid.UUID())
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200 {
		return ioutil.ReadAll(response.Body)
	}
	if response.StatusCode == 500 {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var e500 *Error500
		err = json.Unmarshal(bytes, &e500)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(e500.Errno + ":" + e500.Error)
	}
	return nil, errors.New("99999:HTTP异常(" + response.Status + ")，请检查业务服务 POST " + url)
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
	return Post(ctx, url, body)
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
