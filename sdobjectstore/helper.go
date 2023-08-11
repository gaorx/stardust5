package sdobjectstore

import (
	"github.com/gaorx/stardust5/sderr"
	"strings"
)

const (
	extPlaceholder    = "{ext}"
	md5Placeholder    = "{md5}"
	defaultObjectName = md5Placeholder + "." + extPlaceholder
)

func expandObjectName(src Source, objectName string) (string, error) {
	if objectName == "" {
		objectName = defaultObjectName
	}
	r := objectName
	if strings.Contains(objectName, md5Placeholder) {
		md5Sum := src.Md5()
		if md5Sum == "" {
			return "", sderr.New("md5 error")
		}
		r = strings.Replace(r, md5Placeholder, md5Sum, -1)
	}
	if strings.Contains(objectName, extPlaceholder) {
		ext := strings.TrimPrefix(src.Ext(), ".")
		r = strings.Replace(r, extPlaceholder, ext, -1)
	}
	return r, nil
}
