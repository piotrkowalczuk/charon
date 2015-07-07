package lib

import (
	"crypto/rand"
	"io"
)

// RandomBytesGenerator ...
type RandomBytesGenerator interface {
	GenerateRandomBytes(int) ([]byte, error)
}

// SystemRandomBytesGenerator ...
type SystemRandomBytesGenerator struct {
}

// SystemRandomBytesGenerator creates a random key with the given length in bytes.
func (srbg *SystemRandomBytesGenerator) GenerateRandomBytes(length int) ([]byte, error) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil, err
	}
	return k, nil
}
