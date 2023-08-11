package sdencoding

import (
	"encoding/hex"
	"github.com/gaorx/stardust5/sderr"
)

type hexEncoding struct{}

var (
	Hex Encoding = hexEncoding{}
)

func (e hexEncoding) Encode(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	buff := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(buff, data)
	return buff
}

func (e hexEncoding) EncodeString(data []byte) string {
	return hex.EncodeToString(data)
}

func (e hexEncoding) Decode(encoded []byte) ([]byte, error) {
	if len(encoded) == 0 {
		return []byte{}, nil
	}
	buff := make([]byte, hex.DecodedLen(len(encoded)))
	n, err := hex.Decode(buff, encoded)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return buff[:n], nil
}

func (e hexEncoding) DecodeString(encoded string) ([]byte, error) {
	buff, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return buff, nil
}
