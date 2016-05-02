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
	plain, err := lzjs.DecompressFromBase64(string(b64lzplain))
	if err != nil {
		t.Fatal(err)
	}

	if plain != input {
		log.Fatalf("Could not decrypt raw post data")
	}
}
