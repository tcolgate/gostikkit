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
