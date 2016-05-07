package oss

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// RFC1123Format 获取符合 RFC1123 格式的时间表达。
//
// RFC1123: https://www.ietf.org/rfc/rfc1123.txt
//
// Go 标准库中的 time.Time 的 Format 方法可以接受 time.RFC1123 作为格式化参数，
// 但标准库的做法过于正确😂，并不严格满足 RFC1123 的（错误）要求。
func RFC1123Format(t time.Time) string {
	text := t.UTC().Format(time.RFC1123)
	if strings.HasSuffix(text, "UTC") {
		// 根据 (错误的) RFC1123 文档的规定，对于 UTC 时间应该写成 GMT
		return strings.Replace(text, "UTC", "GMT", -1)
	}
	return text
}

// RFC1123FormatNow 获取符合 RFC1123 格式的当前时间的表达。
//
// 更多信息参考 RFC1123Format 函数文档。
func RFC1123FormatNow() string {
	return RFC1123Format(time.Now())
}

func createResponseError(resp *http.Response) error {
	responseBytes, _ := ioutil.ReadAll(resp.Body)
	message := fmt.Sprintf("[%s]\n%s", resp.Status, string(responseBytes))
	return errors.New(message)
}
