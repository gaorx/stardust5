package sdobjectstore

import (
	"github.com/gaorx/stardust5/sderr"
	"github.com/samber/lo"
	"strings"
)

type Options struct {
	Type           string `json:"type" toml:"type"`
	Root           string `json:"root" toml:"root"`
	Endpoint       string `json:"endpoint" toml:"endpoint"`
	AccessKey      string `json:"access_key" toml:"access_key"`
	AccessSecret   string `json:"access_secret" toml:"access_secret"`
	Bucket         string `json:"bucket" toml:"bucket"`
	Prefix         string `json:"prefix" toml:"prefix"`
	InternalPrefix string `json:"internal_prefix" toml:"internal_prefix"`
}

func New(opts *Options) (Store, error) {
	opts1 := lo.FromPtr(opts)
	switch strings.ToLower(opts1.Type) {
	case "discard":
		return Discard, nil
	case "dir", "directory":
		return Dir{Root: opts1.Root}, nil
	case "aliyun_oss", "aliyun-oss", "aliyunoss":
		aoss, err := NewAliyunOSS(&AliyunOssOptions{
			Endpoint:       opts1.Endpoint,
			AccessKey:      opts1.AccessKey,
			AccessSecret:   opts1.AccessSecret,
			Bucket:         opts1.Bucket,
			Prefix:         opts1.Prefix,
			InternalPrefix: opts1.InternalPrefix,
		})
		if err != nil {
			return nil, sderr.WithStack(err)
		}
		return aoss, nil
	default:
		return nil, sderr.NewWith("illegal type", opts1.Type)
	}
}
