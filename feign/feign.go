package feign

import "context"

// 内网调用接口
type Feign interface {
	Do(ctx context.Context, method string, body interface{}, result interface{}) (err error)
}
