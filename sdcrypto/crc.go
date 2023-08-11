package sdcrypto

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"

	"github.com/gaorx/stardust5/sderr"
)

// CRC32

type CRC32Encrypter struct {
	Encrypter Encrypter
}

func (e *CRC32Encrypter) Encrypt(key, data []byte) ([]byte, error) {
	data = bytes.Clone(data)
	sumBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sumBytes, crc32.ChecksumIEEE(data))
	r, err := e.Encrypter.Encrypt(key, append(data, sumBytes...))
	if err != nil {
		return nil, sderr.Wrap(err, "CRC32 encrypt error")
	}
	return r, nil
}

func (e *CRC32Encrypter) Decrypt(key, encrypted []byte) ([]byte, error) {
	decrypted, err := e.Encrypter.Decrypt(key, encrypted)
	if err != nil {
		return nil, sderr.Wrap(err, "CRC32 decrypt error")
	}
	n := len(decrypted)
	if n < 4 {
		return nil, sderr.New("decrypted is too short")
	}
	data, sumBytes := decrypted[0:n-4], decrypted[n-4:]
	expectant := binary.LittleEndian.Uint32(sumBytes)
	if crc32.ChecksumIEEE(data) != expectant {
		return nil, sderr.New("CRC32 error")
	}
	return data, nil
}
