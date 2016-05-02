package gostikkit

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"

	"github.com/tcolgate/gostikkit/evpkdf"
)

// Appends padding.
func pkcs7Pad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
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
		return nil, fmt.Errorf("invalid padding")
	}
	// check padding
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padlen], nil
}

var opensslmagic = []byte{0x53, 0x61, 0x6c, 0x74, 0x65, 0x64, 0x5f, 0x5f}

func addSalt(ciphertext, salt []byte) []byte {
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

func encrypt(plaintext, password string, salt []byte) []byte {
	keyiv := evpkdf.New(md5.New, []byte(password), salt, len(password)+aes.BlockSize, 1)
	key := keyiv[0:(len(keyiv) - aes.BlockSize)]
	iv := keyiv[(len(keyiv) - aes.BlockSize):len(keyiv)]

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	padded, err := pkcs7Pad([]byte(plaintext), block.BlockSize())

	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return addSalt(ciphertext, salt)
}

func decrypt(ciphertext []byte, password string) []byte {
	ciphertext, salt := stripSalt(ciphertext)

	keyiv := evpkdf.New(md5.New, []byte(password), salt, len(password)+aes.BlockSize, 1)
	key := keyiv[0:(len(keyiv) - aes.BlockSize)]
	iv := keyiv[(len(keyiv) - aes.BlockSize):len(keyiv)]

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

	return plain
}
