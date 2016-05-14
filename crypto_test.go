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
	"encoding/base64"
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

	dec := string(decrypt(encrypt(plaintext, password), password))
	if plaintext != dec {
		t.Fatalf("failed:\nwanted: %v\n got: %v\n", plaintext, dec)
	}
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

func TestCryptNoSalt(t *testing.T) {
	password := key
	plaintext := "exampleplaintext"

	dec := string(decrypt(encrypt(plaintext, password), password))
	if plaintext != dec {
		t.Fatalf("failed:\nwanted: %v\n got: %v\n", plaintext, dec)
	}
}
