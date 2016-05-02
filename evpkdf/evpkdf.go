package evpkdf

import (
	"bytes"
	"hash"
	"io"
)

//New creates a key derivation function that should match the EVP kdf in crypto-js, which, in turn
// should be compatible with openssl's EVP kdf
func New(hash func() hash.Hash, password []byte, salt []byte, keysize int, iterations int) []byte {
	hasher := hash()
	derivedKey := []byte{}
	block := []byte{}

	// Generate key
	for len(derivedKey) < keysize {
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
	return derivedKey[0:keysize]
}
