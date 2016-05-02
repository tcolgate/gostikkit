package gostikkit

import (
	"encoding/base64"
	"log"
	"testing"

	"github.com/tcolgate/gostikkit/lzjs"
)

func TestPostDecrypt(t *testing.T) {
	// input: "hello"
	// http://stikked.luisaranguren.com/view/9e748433#nsXqrTpGtftKqIfu6txaNFNqdyM3A3Cx
	// "U2FsdGVkX1/Z9F8+4gPRLgKf3izs2BYEtqoPMVY/Rw8="
	input := "hello"
	encoded := "U2FsdGVkX1/Z9F8+4gPRLgKf3izs2BYEtqoPMVY/Rw8="
	key := "nsXqrTpGtftKqIfu6txaNFNqdyM3A3Cx"

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
		return
	}

	b64lzplain := decrypt(decoded, key)
	lzplain, err := base64.StdEncoding.DecodeString(string(b64lzplain))

	lzus := []uint16{}
	for i := 0; i < len(lzplain); i += 2 {
		v := uint16(lzplain[i])
		v <<= 8
		v |= uint16(lzplain[i+1])
		lzus = append(lzus, v)
	}
	plain, err := lzjs.Decompress(lzus)
	if err != nil {
		t.Fatal(err)
	}

	if plain != input {
		log.Fatalf("Could not decrypt raw post data")
	}
}
