package sdencoding

import (
	"testing"
)

func TestBase64(t *testing.T) {
	b := []byte("hello 你好")
	testEncodingHelper(t, b, Base64Std)
	testEncodingHelper(t, b, Base64Url)
}
