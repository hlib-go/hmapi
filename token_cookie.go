package hmapi

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"errors"
	"github.com/hlib-go/hmapi/errs"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

/* 客户端 Cookie 存储密文 token字符串  */

const (
	token_secret = "yh2HP7s2"
)

// 生成token
func GenToken(uid string, expiresTime int64) (token string) {
	src := uid + "&" + time.Now().Add(time.Duration(expiresTime)*time.Second).Format(time.RFC3339)
	token, err := DES_ECB_PKCS5_Encode(src, token_secret)
	if err != nil {
		log.Error("GenToken->", err.Error())
	}
	return
}

// 校验token
func VerToken(token string) (uid string, err error) {
	defer func() {
		if err != nil {
			log.Error("VerifyOauthToken Error:" + err.Error())
			err = errs.E99911
		}
	}()
	src, err := DES_ECB_PKCS5_Decode(token, token_secret)
	if err != nil {
		return
	}
	ts := strings.Split(src, "&")
	if len(ts) != 2 {
		err = errs.E99911
		return
	}
	uid = ts[0]
	if uid == "" {
		err = errs.E99911
		return
	}

	// 是否超时
	t, err := time.Parse(time.RFC3339, ts[1])
	if err != nil {
		return
	}
	if time.Now().After(t) {
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
