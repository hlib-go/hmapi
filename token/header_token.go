package token

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/hlib-go/hmapi/errs"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	INVALID_TOKEN = errs.E99911
)

type Object struct {
	Uid     string    `json:"uid"`
	Mobile  string    `json:"mobile"`
	Second  int64     `json:"second"`  // 有效期秒数
	Expires time.Time `json:"expires"` // 到期时间
}

func (t *Object) Json() string {
	tbytes, _ := json.Marshal(t)
	return string(tbytes)
}

func (t *Object) SetExpires(second int64) *Object {
	t.Expires = time.Now().Add(time.Duration(second) * time.Second)
	return t
}

func (t *Object) Gen(secret string) string {
	token, err := des_ecb_pkcs5_encode(t.Json(), secret)
	if err != nil {
		token = "gen-token-error"
		log.Error("Gen Token Error: ", err.Error())
	}
	return token
}

// 生成 Token
func Gen(secret string, t *Object) (tokenVal string) {
	defer func() {
		if e := recover(); e != nil {
			tokenVal = "gen-token-error"
			log.Error(e)
		}
	}()
	tokenVal = t.Gen(secret)
	return
}

// 验证 Token
func Ver(secret, tokenVal string) (t *Object, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
			err = INVALID_TOKEN.NewMsg("解析TOKEN出错")
		}
	}()
	src, err := des_ecb_pkcs5_decode(tokenVal, secret)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(src), &t)
	if err != nil {
		return
	}
	// 验证是否超时
	if time.Now().After(t.Expires) {
		err = INVALID_TOKEN
		return
	}
	if t.Uid == "" {
		err = INVALID_TOKEN
		return
	}
	return
}

// des_ecb_pkcs5_encode
func des_ecb_pkcs5_encode(src, key string) (v string, err error) {
	if len(key) != 8 {
		err = errors.New("Token密钥长度错误")
		return
	}
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return
	}
	bs := block.BlockSize()
	//对明文数据进行补码
	data = pkCS5Padding(data, bs)
	if len(data)%bs != 0 {
		err = errors.New("Need a multiple of the blocksize")
		return
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	v = base64.RawURLEncoding.EncodeToString(out)
	return
}

// des_ecb_pkcs5_decode
func des_ecb_pkcs5_decode(src, key string) (v string, err error) {
	if len(key) != 8 {
		err = errors.New("Token密钥长度错误")
		return
	}
	data, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return
	}
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		err = errors.New("crypto/cipher: input not full blocks")
		return
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = pkCS5UnPadding(out)
	v = string(out)
	return
}

func pkCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
