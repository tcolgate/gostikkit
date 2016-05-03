// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of lzjs.
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

// Package lzjs implements the lz-string data compress / decompression
// routines. It aims for full binary compatibility with lz-string
package lzjs

import (
	"encoding/base64"
	"errors"
	"log"
	"unicode/utf16"
)

type lzData struct {
	val      uint16
	str      []uint16
	position int
	index    int
	empty    bool
}

type lzCtx struct {
	dictionary         map[string]int
	dictionaryToCreate map[string]bool
	c                  uint16
	wc                 []uint16
	w                  []uint16
	enlargeIn          int
	dictSize           int
	numBits            int
	result             []uint16
	data               *lzData
}

func pow2(n int) int {
	return 1 << uint(n)
}

func (data *lzData) writeBit(value uint16) {
	data.val = (data.val << 1) | uint16(value)
	if data.position == 15 {
		data.position = 0
		data.str = append(data.str, data.val)
		data.val = 0
	} else {
		data.position++
	}
}

func (data *lzData) writeBits(numBits int, value uint16) {
	for i := 0; i < numBits; i++ {
		data.writeBit(value & 1)
		value = value >> 1
	}
}

func (ctx *lzCtx) decrementEnlargeIn() {
	ctx.enlargeIn--
	if ctx.enlargeIn == 0 {
		ctx.enlargeIn = pow2(ctx.numBits)
		ctx.numBits++
	}
}

func hashUtf16(i []uint16) string {
	return string(utf16.Decode(i))
}

func (ctx *lzCtx) produceW() {
	if _, ok := ctx.dictionaryToCreate[hashUtf16(ctx.w)]; ok {
		iw := ctx.w[0]
		if iw < 256 {
			ctx.data.writeBits(ctx.numBits, 0)
			ctx.data.writeBits(8, iw)
		} else {
			ctx.data.writeBits(ctx.numBits, 1)
			ctx.data.writeBits(16, iw)
		}
		ctx.decrementEnlargeIn()
		delete(ctx.dictionaryToCreate, hashUtf16(ctx.w))
	} else {
		ctx.data.writeBits(ctx.numBits, uint16(ctx.dictionary[hashUtf16(ctx.w)]))
	}
	ctx.decrementEnlargeIn()
}

// Compress is the basic routine to compress data to uint16, to
// match the utf16 data produced in lz-string
func Compress(uncompressed string) []uint16 {
	ctx := &lzCtx{
		dictionary:         map[string]int{},
		dictionaryToCreate: map[string]bool{},
		c:                  uint16(0),
		wc:                 []uint16{},
		w:                  []uint16{},
		enlargeIn:          2, // Compensate for the first entry which should not count
		dictSize:           3,
		numBits:            2,
		result:             []uint16{},
		data:               &lzData{},
	}

	chars := utf16.Encode([]rune(uncompressed))

	for _, c := range chars {
		ctx.c = c
		hc := hashUtf16([]uint16{ctx.c})
		if _, ok := ctx.dictionary[hc]; !ok {
			ctx.dictionary[hc] = ctx.dictSize
			ctx.dictSize++
			ctx.dictionaryToCreate[hc] = true
		}

		ctx.wc = append(ctx.w, ctx.c)
		hwc := hashUtf16(ctx.wc)
		if _, ok := ctx.dictionary[hwc]; ok {
			ctx.w = make([]uint16, len(ctx.wc))
			copy(ctx.w, ctx.wc)
		} else {
			ctx.produceW()
			// Add wc to the dictionary.
			ctx.dictionary[hwc] = ctx.dictSize
			ctx.dictSize++
			ctx.w = []uint16{ctx.c}
		}
	}

	// Output the code for w.
	if len(ctx.w) != 0 {
		ctx.produceW()
	}

	// Mark the end of the stream
	ctx.data.writeBits(ctx.numBits, 2)

	// Flush the last char
	for ctx.data.val > 0 {
		ctx.data.writeBit(0)
	}

	return ctx.data.str
}

func (data *lzData) readBit() int {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("failed ", r)
		}
	}()
	if data.empty {
		return 0
	}

	res := data.val & uint16(data.position)
	data.position >>= 1
	if data.position == 0 {
		if data.index < len(data.str) {
			data.position = 32768
			data.val = data.str[data.index]
			data.index++
		} else {
			data.empty = true
		}
	}
	//data.val = (data.val << 1);
	if res > 0 {
		return 1
	}
	return 0
}

func (data *lzData) readBits(numBits int) uint16 {
	res := uint16(0)
	maxpower := pow2(numBits)
	var power = 1
	for power != maxpower {
		res |= uint16(data.readBit() * power)
		power <<= 1
	}
	return res
}

// Decompress is the basic routine to decompress the utf16 data produced in lz-string
// into the origin string
func Decompress(compressed []uint16) (string, error) {
	dictionary := [][]uint16{}
	enlargeIn := 4
	numBits := 3
	entry := []uint16{}
	result := []uint16{}
	w := []uint16{}
	var c uint16
	errorCount := 0
	data := &lzData{
		str:      compressed,
		val:      compressed[0],
		position: 32768,
		index:    1,
	}

	for i := 0; i < 3; i++ {
		dictionary = append(dictionary, []uint16{uint16(i)})
	}

	next := data.readBits(2)
	switch next {
	case 0:
		c = data.readBits(8)
		break
	case 1:
		c = data.readBits(16)
		break
	case 2:
		return "", nil
	}
	dictionary = append(dictionary, []uint16{uint16(c)})
	w = []uint16{c}
	result = []uint16{c}

	for {
		c = data.readBits(numBits)

		switch c {
		case 0:
			if errorCount > 10000 {
				return "", errors.New("too many errors")
			}
			errorCount++
			c = data.readBits(8)
			dictionary = append(dictionary, []uint16{uint16(c)})
			c = uint16(len(dictionary)) - 1
			enlargeIn--
			break
		case 1:
			c = data.readBits(16)
			dictionary = append(dictionary, []uint16{uint16(c)})
			c = uint16(len(dictionary)) - 1
			enlargeIn--
			break
		case 2:
			return string(utf16.Decode(result)), nil
		}

		if enlargeIn == 0 {
			enlargeIn = pow2(numBits)
			numBits++
		}

		if int(c) < len(dictionary) {
			entry = dictionary[int(c)]
		} else {
			if c == uint16(len(dictionary)) {
				entry = append(w, w[0])
			} else {
				//return string(utf16.Decode(result)), nil
				return "", errors.New("ran out of dictionary")
			}
		}
		result = append(result, entry...)

		// Add w+entry[0] to the dictionary.
		newdr := make([]uint16, len(w)+1)
		copy(newdr, append(w, entry[0]))
		dictionary = append(dictionary, newdr)
		enlargeIn--

		w = entry

		if enlargeIn == 0 {
			enlargeIn = pow2(numBits)
			numBits++
		}
	}
}

// CompressToBase64 takes a string and returns the
// base64 encoded, compressed data
func CompressToBase64(in string) string {
	comps := Compress(in)

	bs := []byte{}
	for i := 0; i < len(comps); i++ {
		bs = append(bs, byte(comps[i]>>8))
		bs = append(bs, byte(comps[i]&(uint16(255))))
	}
	return base64.StdEncoding.EncodeToString(bs)
}

// DecompressFromBase64 takes a base64 encoded compressed
// blob and returns the origin, uncompressed data
func DecompressFromBase64(in string) (string, error) {
	lzplain, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}

	lzus := []uint16{}
	for i := 0; i < len(lzplain); i += 2 {
		v := uint16(lzplain[i])
		v <<= 8
		v |= uint16(lzplain[i+1])
		lzus = append(lzus, v)
	}
	return Decompress(lzus)
}
