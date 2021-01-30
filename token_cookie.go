package hmapi

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/hlib-go/hmapi/errs"
	"os"
	"time"
)

// 生成密文token字符串的秘钥
var h_token_secret = "yh9HP7s2"

func init() {
	secret := os.Getenv("H_TOKEN_SECRET") // 注意：需所有服务都设置环境变量，否则会导致部分服务token失效
	if secret != "" && len(secret) == 8 {
		h_token_secret = secret
	}
}

type TokenCookie struct {
	Uid     string    `json:"uid"`
	Mobile  string    `json:"mobile"`
	Expires time.Time `json:"expires"` // 到期时间
}

// 生成 Token
func GenToken(uid, mobile string, second int64) (token string) {
	defer func() {
		if e := recover(); e != nil {
			token = "gen-token-error"
		}
	}()
	tc := &TokenCookie{
		Uid:     uid,
		Mobile:  mobile,
		Expires: time.Now().Add(time.Duration(second) * time.Second),
	}
	tbytes, _ := json.Marshal(tc)
	token, err := DES_ECB_PKCS5_Encode(string(tbytes), h_token_secret)
	if err != nil {
		panic(err)
	}
	return
}

// 验证 Token
func VerToken(token string) (t *TokenCookie, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errs.E99911
		}
	}()
	src, err := DES_ECB_PKCS5_Decode(token, h_token_secret)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(src), &t)
	if err != nil {
		return
	}
	// 验证是否超时
	if time.Now().After(t.Expires) {
		t = nil
		err = errs.E99911
		return
	}
	return
}

// DES_ECB_PKCS5_Encode
func DES_ECB_PKCS5_Encode(src, key string) (v string, err error) {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return
	}
	bs := block.BlockSize()
	//对明文数据进行补码
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		err = errors.New("Need a multiple of the blocksize")
		return
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//对明文按照blocksize进行分块加密
		//必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	v = base64.RawURLEncoding.EncodeToString(out)
	return
}

// DES_ECB_PKCS5_Decode
func DES_ECB_PKCS5_Decode(src, key string) (v string, err error) {
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
	out = PKCS5UnPadding(out)
	v = string(out)
	return
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
