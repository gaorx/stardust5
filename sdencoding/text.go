package sdencoding

import (
	"github.com/gaorx/stardust5/sderr"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type TextEncoding interface {
	Encode(s string) ([]byte, error)
	Decode(encoded []byte) (string, error)
}

var (
	GBK     TextEncoding = textEncoding{simplifiedchinese.GBK}
	GB2312  TextEncoding = textEncoding{simplifiedchinese.HZGB2312}
	GB18030 TextEncoding = textEncoding{simplifiedchinese.GB18030}
)

type textEncoding struct {
	encoding encoding.Encoding
}

func (e textEncoding) Encode(s string) ([]byte, error) {
	b, err := e.encoding.NewEncoder().Bytes([]byte(s))
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return b, nil
}

func (e textEncoding) Decode(encoded []byte) (string, error) {
	if encoded == nil {
		return "", sderr.New("encoded is nil")
	}
	b, err := e.encoding.NewDecoder().Bytes(encoded)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return string(b), nil
}
