package sdencoding

type Encoding interface {
	Encode(data []byte) []byte
	EncodeString(data []byte) string
	Decode(encoded []byte) ([]byte, error)
	DecodeString(encoded string) ([]byte, error)
}
