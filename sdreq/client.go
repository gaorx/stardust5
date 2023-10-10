package sdreq

import (
	"crypto/tls"
	"github.com/gaorx/stardust5/sdtime"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"net/http"
)

type Options struct {
	// common
	BaseUrl      string
	UserAgent    string
	Headers      map[string]string
	QueryParams  map[string]string
	PathParams   map[string]string
	AuthUsername string
	AuthPassword string
	Cookies      []*http.Cookie

	// cert
	Certs        []tls.Certificate
	CertFilename string
	KeyFilename  string

	// log/dump
	Log      bool
	DebugLog bool
	DumpAll  bool

	// timeout / retry
	TimeoutMS       int64 // 超时时间，单位毫秒
	RetryIntervalMS int64
	RetryCount      int
}

func New(opts *Options) *req.Client {
	opts1 := lo.FromPtr(opts)
	c := req.C()
	if opts1.BaseUrl != "" {
		c.SetBaseURL(opts1.BaseUrl)
	}
	if opts1.UserAgent != "" {
		c.SetUserAgent(opts1.UserAgent)
	}
	if opts1.TimeoutMS > 0 {
		c.SetTimeout(sdtime.Milliseconds(opts1.TimeoutMS))
	}
	if len(opts1.Headers) > 0 {
		c.SetCommonHeaders(opts1.Headers)
	}
	if len(opts1.QueryParams) > 0 {
		c.SetCommonQueryParams(opts1.QueryParams)
	}
	if len(opts1.PathParams) > 0 {
		c.SetCommonPathParams(opts1.PathParams)
	}
	if opts1.AuthUsername != "" {
		c.SetCommonBasicAuth(opts1.AuthUsername, opts1.AuthPassword)
	}
	if len(opts1.Cookies) > 0 {
		c.SetCommonCookies(opts1.Cookies...)
	}
	if len(opts1.Certs) > 0 {
		c.SetCerts(opts1.Certs...)
	}
	if opts1.CertFilename != "" && opts1.KeyFilename != "" {
		c.SetCertFromFile(opts1.CertFilename, opts1.KeyFilename)
	}
	if opts1.Log {
		c.SetLogger(DefaultLogger)
	}
	if opts1.DebugLog {
		c.EnableDebugLog()
	}
	if opts1.DumpAll {
		c.EnableDumpAll()
	}
	if opts1.RetryIntervalMS > 0 {
		c.SetCommonRetryFixedInterval(sdtime.Milliseconds(opts1.RetryIntervalMS))
	}
	if opts1.RetryCount > 0 {
		c.SetCommonRetryCount(opts1.RetryCount)
	}
	return c
}
