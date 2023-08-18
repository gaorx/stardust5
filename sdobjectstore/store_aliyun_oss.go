package sdobjectstore

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
	"io"
	"strings"
)

type AliyunOssOptions struct {
	Endpoint       string `json:"endpoint" toml:"endpoint"`
	AccessKey      string `json:"access_key" toml:"access_key"`
	AccessSecret   string `json:"access_secret" toml:"access_secret"`
	Bucket         string `json:"bucket" toml:"bucket"`
	Prefix         string `json:"prefix" toml:"prefix"`
	InternalPrefix string `json:"internal_prefix" toml:"internal_prefix"`
}

type aliyunOss struct {
	client         *oss.Client
	bucket         string
	prefix         string
	internalPrefix string
}

func NewAliyunOSS(opts *AliyunOssOptions) (Interface, error) {
	opts1 := lo.FromPtr(opts)
	if opts1.Bucket == "" {
		return nil, sderr.New("no bucket")
	}
	if opts1.Prefix == "" {
		return nil, sderr.New("no prefix")
	}
	if opts1.InternalPrefix == "" {
		return nil, sderr.New("no internal prefix")
	}
	client, err := oss.New(opts1.Endpoint, opts1.AccessKey, opts1.AccessSecret)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return &aliyunOss{
		client:         client,
		bucket:         opts1.Bucket,
		prefix:         opts1.Prefix,
		internalPrefix: opts1.InternalPrefix,
	}, nil
}

func (aoss *aliyunOss) Store(src Source, objectName string) (*Target, error) {
	if src == nil {
		return nil, sderr.New("nil source")
	}

	// 展开文件名称
	expandedObjectName, err := expandObjectName(src, objectName)
	if err != nil {
		return nil, err
	}

	// 针对特殊类型content-type，修正掉content-disposition
	ct, cd := src.ContentType(), "attachment"
	if strings.HasPrefix(ct, "text/") ||
		strings.HasPrefix(ct, "image/") ||
		strings.HasPrefix(ct, "video/") ||
		strings.HasPrefix(ct, "audio/") ||
		strings.Contains(ct, "json") ||
		strings.Contains(ct, "javascript") ||
		strings.Contains(ct, "ecmascript") {
		cd = "inline" // 阿里云实际上会忽略这个inline，强制改为attachment
	}

	// 存储
	b, err := aoss.client.Bucket(aoss.bucket)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	r, err := src.Open()
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	defer func() {
		if rc, ok := r.(io.ReadCloser); ok {
			_ = rc.Close()
		}
	}()
	err = b.PutObject(expandedObjectName, r, oss.ContentType(src.ContentType()), oss.ContentDisposition(cd))
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	// 返回信息
	return &Target{
		Typ:            HttpTarget,
		Prefix:         aoss.prefix,
		InternalPrefix: aoss.internalPrefix,
		Path:           expandedObjectName,
	}, nil
}
