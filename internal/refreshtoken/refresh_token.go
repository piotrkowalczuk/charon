package refreshtoken

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/sha3"
)

// Random generate refresh token with given key and generated hash of length 16.
func Random() (string, error) {
	buf, err := generateRandomBytes(64)
	if err != nil {
		return "", err
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	hash := make([]byte, 32)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum128(hash, buf)
	hash2 := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(hash2, hash)
	return string(hash2), nil
}

func generateRandomBytes(length int) ([]byte, error) {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil, err
	}
	return k, nil
}
