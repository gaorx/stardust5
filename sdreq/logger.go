package sdreq

import (
	"github.com/gaorx/stardust5/sdslog"
	"github.com/imroc/req/v3"
)

var DefaultLogger req.Logger = logger{}

type logger struct {
}

func (_ logger) Errorf(format string, v ...any) {
	sdslog.Errorf(format, v...)
}

func (_ logger) Warnf(format string, v ...any) {
	sdslog.Warnf(format, v...)
}

func (_ logger) Debugf(format string, v ...any) {
	sdslog.Debugf(format, v...)
}
