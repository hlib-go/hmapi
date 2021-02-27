# hmapi

## 介绍

hmapi

## 接口规范

- 请求地址：https://xxx.xx/api/gateway?requestId=xxxx
- 请求方式：POST 表单请求
- 请求报文：appid=xxx&method=hm.xxx.xxx&biz=xxxxx&encrypt=DES_MD5
- 响应报文：JSON

## 报文加解密

### 方式1：DES_MD5

对报文进行DES加密，对加密结果做MD5签名。 拼接签名结果与加密结果得到最终密文字符串

### 方式2：AES_SHA256

AES加密后，计算SHA256签名，签名拼接密文得到最终密文

### 方式3：RSA

RSA加密

### 方式4：AES_RSA

生成随机密钥，使用AES加密报文，使用RSA加密随机密钥。 密钥密文拼接数据密文得最终结果