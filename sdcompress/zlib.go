package sdcompress

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"

	"github.com/gaorx/stardust5/sderr"
)

type ZlibLevel int

const (
	ZlibNoCompression      = ZlibLevel(zlib.NoCompression)
	ZlibBestSpeed          = ZlibLevel(zlib.BestSpeed)
	ZlibBestCompression    = ZlibLevel(zlib.BestCompression)
	ZlibDefaultCompression = ZlibLevel(zlib.DefaultCompression)
	ZlibHuffmanOnly        = ZlibLevel(zlib.HuffmanOnly)
)

var (
	ZlibAllLevels = []ZlibLevel{
		ZlibNoCompression,
		ZlibBestSpeed,
		ZlibBestCompression,
		ZlibDefaultCompression,
		ZlibHuffmanOnly,
	}
)

func Zip(data []byte, level ZlibLevel) ([]byte, error) {
	if data == nil {
		return nil, sderr.New("zip nil data")
	}
	buff := new(bytes.Buffer)
	w, err := zlib.NewWriterLevel(buff, int(level))
	if err != nil {
		return nil, sderr.WrapWith(err, "zip make writer error", level)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, sderr.Wrap(err, "zip write error")
	}
	err = w.Close()
	if err != nil {
		return nil, sderr.Wrap(err, "zip close error")
	}
	return buff.Bytes(), nil
}

func Unzip(data []byte) ([]byte, error) {
	if data == nil {
		return nil, sderr.New("unzip nil data")
	}
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, sderr.Wrap(err, "unzip make reader error")
	}
	defer func() { _ = r.Close() }()

	to, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, sderr.Wrap(err, "unzip read error")
	}
	return to, nil
}
