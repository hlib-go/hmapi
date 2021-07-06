# hmapi

## 介绍

hmapi

## 接口规范

- 请求地址：https://xxx.xx/api/gateway?requestId=xxxx
- 请求方式：POST 表单请求
- 请求报文：appid=xxx&method=hm.xxx.xxx&biz=xxxxx&encrypt=v2
- 响应报文：JSON

## 报文加解密

### v1：DES_MD5

对报文进行DES加密，对加密结果做MD5签名。 拼接签名结果与加密结果得到最终密文字符串

### v2：AES_RSA

生成随机密钥，使用AES加密报文，使用RSA加密随机密钥。 密钥密文拼接数据密文得最终结果

### v3：AES_SHA1_RSA
