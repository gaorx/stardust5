package sdtime

import (
	"time"
)

// time -> unix

func ToUnixS(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}

func ToUnixMS(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / int64(time.Millisecond)
}

// unix -> time

func FromUnixS(s int64) time.Time {
	if s == 0 {
		return time.Time{}
	}
	return time.Unix(s, 0)
}

func FromUnixMS(ms int64) time.Time {
	if ms == 0 {
		return time.Time{}
	}
	nanos := ms * 1000000
	return time.Unix(0, nanos)
}

// now in unix

func NowUnixS() int64 {
	return ToUnixS(time.Now())
}

func NowUnixMS() int64 {
	return ToUnixMS(time.Now())
}
