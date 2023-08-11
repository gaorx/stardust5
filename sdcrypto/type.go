package sdcrypto

type Encrypter interface {
	Encrypt(key, data []byte) ([]byte, error)
	Decrypt(key, crypted []byte) ([]byte, error)
}

// Encode

type EncrypterFunc struct {
	Encrypter func(key, data []byte) ([]byte, error)
	Decrypter func(key, crypted []byte) ([]byte, error)
}

func (e *EncrypterFunc) Encrypt(key, data []byte) ([]byte, error) {
	return e.Encrypter(key, data)
}

func (e *EncrypterFunc) Decrypt(key, crypted []byte) ([]byte, error) {
	return e.Decrypter(key, crypted)
}
