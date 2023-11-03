package sdurl

import (
	"fmt"
	"github.com/gaorx/stardust5/sdstrings"
	"strings"
)

func SplitHostPort(hostPort string) (host, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		host, port = host[:colon], host[colon+1:]
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

func CompleteHttp(urlStr string, defaultScheme string) string {
	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		return urlStr
	}
	return fmt.Sprintf(
		"%s://%s",
		sdstrings.TrimSuffixes(defaultScheme, "//", ":"),
		sdstrings.TrimPrefixes(urlStr, "//", "/"),
	)
}

func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}
