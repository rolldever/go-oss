package oss

import (
	"fmt"
	"strings"
)

var (
	// OSSServerDomainName 是OSS 服务器的域名。
	// 包的使用者可以直接修改这个值。
	OSSServerDomainName = "aliyuncs.com"
)

// Bucket 抽象了 OSS 中的一个 Bucket.
type Bucket struct {
	Name       string
	Datacenter string
	*Auth
}

// NewBucket 创建一个 Bucket 对象，必须提供 bucket 的名字和数据中心，而 auth 可以为 nil.
// 如果 auth 为 nil，相当于调用 NewPublicBucket，即创建一个公共 Bucket 对象。
func NewBucket(name, datacenter string, auth *Auth) *Bucket {
	return &Bucket{
		Name:       name,
		Datacenter: datacenter,
		Auth:       auth,
	}
}

// NewPublicBucket 通过给定的 Bucket 名称和数据中心创建一个公共的 Bucket 对象。
func NewPublicBucket(name, datacenter string) *Bucket {
	return NewBucket(name, datacenter, nil)
}

// GetObjectURL 获取给定的 Object 的 URL.
func (b *Bucket) GetObjectURL(objectName string) string {
	return b.getObjectURL(objectName, false)
}

// GetObjectURLHTTPS 获取给定的 Object 的 URL (使用 HTTPS 协议).
func (b *Bucket) GetObjectURLHTTPS(objectName string) string {
	return b.getObjectURL(objectName, true)
}

func (b *Bucket) getObjectURL(objectName string, usingHTTPS bool) string {
	if strings.HasPrefix(objectName, "/") {
		objectName = objectName[1:]
	}

	var protocol string
	if usingHTTPS {
		protocol = "https"
	} else {
		protocol = "http"
	}
	return fmt.Sprintf("%s://%s.%s.%s/%s",
		protocol, b.Name, b.Datacenter, OSSServerDomainName, objectName)
}
