package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
)

// Auth 封装用户认证信息。
type Auth struct {
	PublicKey string
	SecretKey string
}

// NewAuth 通过给定的 Public Key 和 Secret Key 创建一个 Auth.
func NewAuth(publicKey, secretKey string) *Auth {
	return &Auth{
		PublicKey: publicKey,
		SecretKey: secretKey,
	}
}

// SignRequest 为 HTTP 请求签名。
func (auth *Auth) SignRequest(req *http.Request, bucketName, objectName string) {
	buff := new(bytes.Buffer)
	buff.WriteString(req.Method)
	buff.WriteRune('\n')

	buff.WriteString(req.Header.Get("Content-MD5"))
	buff.WriteRune('\n')

	buff.WriteString(req.Header.Get("Content-Type"))
	buff.WriteRune('\n')

	buff.WriteString(req.Header.Get("Date"))
	buff.WriteRune('\n')

	// CanonicalizedResource
	buff.WriteString("/")
	buff.WriteString(bucketName)
	buff.WriteString("/")
	buff.WriteString(objectName)

	mac := hmac.New(sha1.New, []byte(auth.SecretKey))
	mac.Write(buff.Bytes())
	authString := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	authString = fmt.Sprintf("OSS %s:%s", auth.PublicKey, authString)
	req.Header.Set("Authorization", authString)
}
