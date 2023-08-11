package sdbytes

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

type Slice []byte

func (s Slice) Clone() Slice {
	return bytes.Clone(s)
}

func (s Slice) String() string {
	if s == nil {
		return "bytes(0) nil"
	}
	l := len(s)
	switch len(s) {
	case 0:
		return "bytes(0) []"
	case 1:
		return fmt.Sprintf("bytes(1) [%s]", hexByte(s[0]))
	case 2:
		return fmt.Sprintf("bytes(2) [%s %s]", hexByte(s[0]), hexByte(s[1]))
	case 3:
		return fmt.Sprintf("bytes(3) [%s %s %s]", hexByte(s[0]), hexByte(s[1]), hexByte(s[2]))
	case 4:
		return fmt.Sprintf("bytes(4) [%s %s %s %s]", hexByte(s[0]), hexByte(s[1]), hexByte(s[2]), hexByte(s[3]))
	default:
		return fmt.Sprintf("bytes(%d) [%s %s ... %s %s]", l, hexByte(s[0]), hexByte(s[1]), hexByte(s[l-2]), hexByte(s[l-1]))
	}
}

func (s Slice) HexL() string {
	if len(s) <= 0 {
		return ""
	}
	return hex.EncodeToString(s)
}

func (s Slice) HexU() string {
	if len(s) <= 0 {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(s))
}

func (s Slice) Base64Std() string {
	if len(s) <= 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(s)
}

func (s Slice) Base64Url() string {
	if len(s) <= 0 {
		return ""
	}
	return base64.URLEncoding.EncodeToString(s)
}

const hextable = "0123456789abcdef"

func hexByte(b byte) string {
	buff := [2]byte{}
	buff[0] = hextable[b>>4]
	buff[1] = hextable[b&0x0f]
	return string(buff[:])
}
