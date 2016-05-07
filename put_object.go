package oss

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"net/http"
)

// PutObjectSettings 用来封装 PUT 操作的个性化配置。
type PutObjectSettings struct {
	ContentType     string
	ContentEncoding string
	Disposition     string
	CalcMD5         bool
}

// NewDefaultPutObjectSettings 创建一个默认的 PutObjectSettings 对象。
// 如果调用 PutObject 方法上传对象，就会使用这个默认的 PutObjectSettings
func NewDefaultPutObjectSettings() *PutObjectSettings {
	return &PutObjectSettings{
		ContentType:     ContentTypeBinary,
		ContentEncoding: "",
		Disposition:     "",
		CalcMD5:         true,
	}
}

// PutObject 执行 OSS Put Object 操作
//
// PutObject 在执行完全成功后才返回 nil，
// 即除非 HTTP 请求成功，并且服务器返回 200 状态码。
func (bucket *Bucket) PutObject(payload []byte, objectName string) error {
	return bucket.PutObjectSettings(payload, objectName, nil)
}

// PutObjectSettings 执行 OSS Put Object 操作，同时提供一组 sessings.
//
// PutObject 在执行完全成功后才返回 nil，
// 即除非 HTTP 请求成功，并且服务器返回 200 状态码。
func (bucket *Bucket) PutObjectSettings(payload []byte, objectName string, settings *PutObjectSettings) error {
	// 如果没有提供要 put 的数据，作为错误处理
	if payload == nil || len(payload) == 0 {
		return errors.New("No content to put")
	}

	// 如果没有提供 settings，则创建一个默认的 settings
	if settings == nil {
		settings = DefaultPutObjectSettings()
	}

	// 如果 settings 中没有提供任何 ContentType，
	// 使用 Binary 形式
	if len(settings.ContentType) == 0 {
		settings.ContentType = ContentTypeBinary
	}

	// 构造 HTTP 请求
	url := bucket.GetObjectURL(objectName)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", settings.ContentType)
	req.ContentLength = int64(len(payload))
	req.Header.Set("Date", RFC1123FormatNow())

	if len(settings.ContentEncoding) > 0 {
		req.Header.Set("Content-Encoding", settings.ContentEncoding)
	}
	if len(settings.Disposition) > 0 {
		req.Header.Set("Content-Disposition", settings.Disposition)
	}

	// 如果 settings 中要求计算 MD5
	if settings.CalcMD5 {
		h := md5.New()
		h.Write(payload)
		md5Bytes := h.Sum(nil)
		req.Header.Set("Content-MD5", base64.StdEncoding.EncodeToString(md5Bytes))
	}

	// 如果提供了 Auth，为请求签名
	if bucket.Auth != nil {
		bucket.Auth.SignRequest(req, bucket.Name, objectName)
	}

	// 发送 HTTP 请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return createResponseError(resp)
	}

	return nil
}
