package sdurl

import (
	"net/url"

	"github.com/gaorx/stardust5/sderr"
)

type Modifier func(u *url.URL)

func Modify(rawUrl string, modifiers ...Modifier) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", sderr.WrapWith(err, "parse url for modify error", rawUrl)
	}
	for _, f := range modifiers {
		if f != nil {
			f(u)
		}
	}
	return u.String(), nil
}

func ModifyDef(rawUrl string, def string, modifiers ...Modifier) string {
	r, err := Modify(rawUrl, modifiers...)
	if err != nil {
		return def
	}
	return r
}

func SetQuery(k, v string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		q.Set(k, v)
		u.RawQuery = q.Encode()
	}
}

func SetQueries(m map[string]string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		for k, v := range m {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
}

func DeleteQuery(k string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		q.Del(k)
		u.RawQuery = q.Encode()
	}
}

func DeleteQueries(keys ...string) Modifier {
	return func(u *url.URL) {
		q := u.Query()
		for _, k := range keys {
			q.Del(k)
		}
		u.RawQuery = q.Encode()
	}
}

func SetPath(path string) Modifier {
	return func(u *url.URL) {
		u.Path = path
	}
}

func SetHost(host string) Modifier {
	return func(u *url.URL) {
		u.Host = host
	}
}

func SetHostname(hostname string) Modifier {
	return func(u *url.URL) {
		_, port := SplitHostPort(u.Host)
		if port != "" {
			u.Host = hostname + ":" + port
		} else {
			u.Host = hostname
		}
	}
}

func SetPort(port string) Modifier {
	return func(u *url.URL) {
		hostname, _ := SplitHostPort(u.Host)
		if port != "" {
			u.Host = hostname + ":" + port
		} else {
			u.Host = hostname
		}
	}
}
