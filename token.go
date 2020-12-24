package hmapi

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hlib-go/hredis"
)

// Token操作接口定义
type Token interface {
	Gen(uid string, validSecond int64) (token string, err error)
	Get(token string) (uid string)
	Verify(token, uid string) error
}

// hb/hc   接口token设置与校验
// 使用Redis存储Token，初始加载时创建
func NewToken(client *redis.Client) Token {
	return &token{
		Kv: &hredis.Kv{
			KeyPre: "token:", //用户会话Token ,key为token:uuid   value为用户id
			Client: client,
		},
	}
}

type token struct {
	Kv *hredis.Kv
}

//生成token
func (t *token) Gen(uid string, validSecond int64) (token string, err error) {
	token = Rand32()
	err = t.Kv.Set(context.Background(), token, uid, validSecond)
	return
}

// 根据token获取用户id
func (t *token) Get(token string) (uid string) {
	return t.Kv.Get(context.Background(), token)
}

// 验证token,判断token对应的内容是否等于用户编号
func (t *token) Verify(token, uid string) error {
	v := t.Kv.Get(context.Background(), token)
	if v == "" || v != uid {
		return E99911
	}
	return nil
}

// Rand32 使用crypto/rand 随机赋值byte数组， 然后md5返回32位十六进制字符串
func Rand32() string {
	var b = make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum(b))
}
