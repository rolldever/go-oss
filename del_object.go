package oss

import (
	"net/http"
)

// DeleteObject 执行 OSS 的 DELETE 操作。
func (bucket *Bucket) DeleteObject(objectName string) error {
	req, err := http.NewRequest("DELETE", bucket.getObjectURL(objectName, true), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Date", RFC1123FormatNow())
	req.ContentLength = 0

	if bucket.Auth != nil {
		bucket.Auth.SignRequest(req, bucket.Name, objectName)
	}

	// 发送 HTTP 请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// DELETE 请求的正确响应是 204
	if resp.StatusCode != http.StatusNoContent {
		return createResponseError(resp)
	}
	return nil
}
