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

package evpkdf

import (
	"bytes"
	"crypto/md5"
	"testing"
)

func TestVector(t *testing.T) {
	expect := []byte{0xfd, 0xbd, 0xf3, 0x41, 0x9f, 0xff, 0x98, 0xbd, 0xb0, 0x24, 0x13, 0x90, 0xf6, 0x2a, 0x9d, 0xb3, 0x5f, 0x4a, 0xba, 0x29, 0xd7, 0x75, 0x66, 0x37, 0x79, 0x97, 0x31, 0x4e, 0xbf, 0xc7, 0x9, 0xf2, 0xb, 0x5c, 0xa7, 0xb1, 0x8, 0x1f, 0x94, 0xb1, 0xac, 0x12, 0xe3, 0xc8, 0xba, 0x87, 0xd0, 0x5a}
	k := New(md5.New, []byte("password"), []byte("saltsalt"), (256+128)/8, 1)
	if bytes.Compare(expect, k) != 0 {
		t.Fatalf("failed\n expected: %v\n got: %v\n", expect, k)
	}
}
