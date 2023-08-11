package sdhash

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"github.com/gaorx/stardust5/sdbytes"
)

func Sha1(data []byte) sdbytes.Slice {
	sum := sha1.Sum(data)
	return sum[:]
}

func Sha256(data []byte) sdbytes.Slice {
	sum := sha256.Sum256(data)
	return sum[:]
}

func Sha512(data []byte) sdbytes.Slice {
	sum := sha512.Sum512(data)
	return sum[:]
}

func HmacSha1(data, key []byte) sdbytes.Slice {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func HmacSha256(data, key []byte) sdbytes.Slice {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func HmacSha512(data, key []byte) sdbytes.Slice {
	mac := hmac.New(sha512.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func ValidHmacSha1(data, key, expected []byte) bool {
	actual := HmacSha1(data, key)
	return hmac.Equal(actual, expected)
}

func ValidHmacSha256(data, key, expected []byte) bool {
	actual := HmacSha256(data, key)
	return hmac.Equal(actual, expected)
}

func ValidHmacSha512(data, key, expected []byte) bool {
	actual := HmacSha512(data, key)
	return hmac.Equal(actual, expected)
}
