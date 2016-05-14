// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of gostikkit.
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

package gostikkit

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"fmt"

	"github.com/tcolgate/gostikkit/evpkdf"
)

// Appends padding.
func pkcs7Pad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := uint8(1)
	for ((len(data) + int(padlen)) % blocklen) != 0 {
		padlen++
	}

	if int(padlen) > blocklen {
		panic(fmt.Sprintf("generated invalid padding length %v for block length %v", padlen, blocklen))
	}
	pad := bytes.Repeat([]byte{byte(padlen)}, int(padlen))
	return append(data, pad...), nil
}

// Returns slice of the original data without padding.
func pkcs7Unpad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}

	padlen := int(data[len(data)-1])
	if padlen > blocklen || padlen == 0 {
		// Not padded
		return data, nil
	}
	// check padding
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return data, nil
		}
	}

	return data[:len(data)-padlen], nil
}

var opensslmagic = []byte{0x53, 0x61, 0x6c, 0x74, 0x65, 0x64, 0x5f, 0x5f}

func addSalt(ciphertext, salt []byte) []byte {
	if len(salt) == 0 {
		return ciphertext
	}
	return append(append(opensslmagic, salt...), ciphertext...)
}

func stripSalt(ciphertext []byte) ([]byte, []byte) {
	if bytes.Compare(opensslmagic, ciphertext[0:len(opensslmagic)]) == 0 {
		return ciphertext[len(opensslmagic)+8 : len(ciphertext)],
			ciphertext[len(opensslmagic) : len(opensslmagic)+8]
	} else {
		return ciphertext, nil
	}
}

func encrypt(plaintext, password string) []byte {
	salt := genChars(8)

	keylen := 32
	key := make([]byte, keylen)
	ivlen := aes.BlockSize
	iv := make([]byte, ivlen)

	keymat := evpkdf.New(md5.New, []byte(password), salt, keylen+ivlen, 1)
	keymatbuf := bytes.NewReader(keymat)

	n, err := keymatbuf.Read(key)
	if n != keylen || err != nil {
		panic("keymaterial was short reading key")
	}

	n, err = keymatbuf.Read(iv)
	if n != ivlen || err != nil {
		panic("keymaterial was short reading iv")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	padded, err := pkcs7Pad([]byte(plaintext), block.BlockSize())
	if err != nil {
		panic("padding blew up, " + err.Error())
	}

	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return addSalt(ciphertext, salt)
}

func decrypt(ciphertext []byte, password string) []byte {
	ciphertext, salt := stripSalt(ciphertext)

	keylen := 32
	key := make([]byte, keylen)
	ivlen := aes.BlockSize
	iv := make([]byte, ivlen)

	keymat := evpkdf.New(md5.New, []byte(password), salt, keylen+ivlen, 1)
	keymatbuf := bytes.NewReader(keymat)

	n, err := keymatbuf.Read(key)
	if n != keylen || err != nil {
		panic("keymaterial was short reading key")
	}

	n, err = keymatbuf.Read(iv)
	if n != ivlen || err != nil {
		panic("keymaterial was short reading iv")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	plain, err := pkcs7Unpad(ciphertext, block.BlockSize())
	if err != nil {
		panic("padding blew up, " + err.Error())
	}

	return plain
}

var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func genChars(n int) []byte {
	out := make([]byte, n)
	rs := make([]byte, n)
	rand.Read(rs)
	for i := 0; i < n; i++ {
		out[i] = chars[uint(rs[i])%uint(len(chars))]
	}
	return out
}
