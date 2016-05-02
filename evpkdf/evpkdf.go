package evpkdf

import (
	"bytes"
	"hash"
	"io"
)

//crypto-js uses openssl's evpKDF by default
func New(hash func() hash.Hash, password []byte, salt []byte, keysize int, iterations int) io.Reader {
	hasher := hash()
	derivedKey := []byte{}
	block := []byte{}

	derivedKeyWords := []uint16{}

	// Generate key
	for len(derivedKeyWords) < keySize {
		if len(block) != 0 {
			io.Copy(hasher, bytes.NewBuffer(block))
		}
		io.Copy(hasher, bytes.NewBuffer(password))
		io.Copy(hasher, bytes.NewBuffer(salt))
		block = hasher.Sum(nil)
		hasher.Reset()

		// Iterations
		for i := 1; i < iterations; i++ {
			io.Copy(hasher, bytes.NewBuffer(block))
			block = hasher.Sum(nil)
			hasher.Reset()
		}

		derivedKey = append(derivedKey, block...)
	}
	derivedKey.sigBytes = keySize * 4

	return bytes.NewBuffer(derivedKey)
}
