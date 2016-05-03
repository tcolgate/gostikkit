// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of evpkdf.
//
// radia is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// radia is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with radia.  If not, see <http://www.gnu.org/licenses/>.

// Package evpkdf implements OpenSSL EVP Key derivation function, aiming to
// be compatible with crypto-js default behaviour.
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
