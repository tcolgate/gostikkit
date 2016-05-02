package gostikkit

import (
	"encoding/base64"
	"log"
	"testing"
)

var plain = `AAA`
var key = `23XIKa66q3fM9hfxxPvIaKaSK584kolA`
var b64 = `U2FsdGVkX19843ZgwE2B9A88goyvERASPauARUgY9HI=`
var lzplain = ``

func TestCrypt1(t *testing.T) {
	//password := "example"
	password := key
	plaintext := "exampleplaintext"

	log.Println(string(decrypt(encrypt(plaintext, password, []byte("saltsalt")), password)))
}

func TestDecrypt1(t *testing.T) {
	//	CryptoJS.AES.decrypt("U2FsdGVkX1+SPCDRNc9kniYWDtYoGS3h/M6xRv2RuYk=","12345678901234567890123456789012").toString(CryptoJS.enc.Utf8)
	ciphertext, err := base64.StdEncoding.DecodeString("U2FsdGVkX1+SPCDRNc9kniYWDtYoGS3h/M6xRv2RuYk=")
	if err != nil {
		t.Fatalf("%v", err)
	}

	password := "12345678901234567890123456789012"
	expected := "hello"

	plaintext := decrypt(ciphertext, password)

	if string(plaintext) != expected {
		t.Fatalf("Wrong decryption result:\n got %v\n expected %v\n", string(plaintext), expected)
	}
}
