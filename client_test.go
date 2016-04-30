package gostikkit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"testing"
)

var plain = `AAA`
var key = `23XIKa66q3fM9hfxxPvIaKaSK584kolA`
var b64 = `U2FsdGVkX19843ZgwE2B9A88goyvERASPauARUgY9HI=`
var lzplain = ``

func TestCrypt1(t *testing.T) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plain))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plain))

	return encodeBase64(ciphertext)
}
