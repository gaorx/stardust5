package sdencoding

import (
	"encoding/base64"
	"github.com/gaorx/stardust5/sderr"
)

type base64Encoding struct {
	encoding *base64.Encoding
}

var (
	Base64Std Encoding = base64Encoding{base64.StdEncoding}
	Base64Url Encoding = base64Encoding{base64.URLEncoding}
)

func (e base64Encoding) Encode(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	enc := e.encoding
	buff := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(buff, data)
	return buff
}

func (e base64Encoding) EncodeString(data []byte) string {
	return e.encoding.EncodeToString(data)
}

func (e base64Encoding) Decode(encoded []byte) ([]byte, error) {
	if len(encoded) == 0 {
		return []byte{}, nil
	}
	enc := e.encoding
	buff := make([]byte, enc.DecodedLen(len(encoded)))
	n, err := enc.Decode(buff, encoded)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return buff[:n], nil
}

func (e base64Encoding) DecodeString(encoded string) ([]byte, error) {
	buff, err := e.encoding.DecodeString(encoded)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return buff, nil
}
