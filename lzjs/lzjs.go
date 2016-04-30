package lzjs

import (
	"math"
	"strings"
)

type lzData struct {
	str      []uint16
	val      uint16
	position int
	index    int
}

type lzCtx struct {
	dictionary         map[string]int
	dictionaryToCreate map[string]bool
	c                  rune
	wc                 string
	w                  string
	enlargeIn          int
	dictSize           int
	numBits            int
	result             []uint16
	data               *lzData
}

func (data *lzData) writeBit(value int) {
	data.val = (data.val << 1) | uint16(value)
	if data.position == 15 {
		data.position = 0
		data.str = append(data.str, data.val)
		data.val = 0
	} else {
		data.position++
	}
}

func (data *lzData) writeBits(numBits int, value int) {
	for i := 0; i < numBits; i++ {
		data.writeBit(value & 1)
		value = value >> 1
	}
}

func (ctx *lzCtx) decrementEnlargeIn() {
	ctx.enlargeIn--
	if ctx.enlargeIn == 0 {
		ctx.enlargeIn = int(math.Pow(2, float64(ctx.numBits)))
		ctx.numBits++
	}
}

func (ctx *lzCtx) produceW() {
	if _, ok := ctx.dictionaryToCreate[ctx.w]; ok {
		iw := strings.IndexRune(ctx.w, 0)
		if iw < 256 {
			ctx.data.writeBits(ctx.numBits, 0)
			ctx.data.writeBits(8, iw)
		} else {
			ctx.data.writeBits(ctx.numBits, 1)
			ctx.data.writeBits(16, iw)
		}
		ctx.decrementEnlargeIn()
		delete(ctx.dictionaryToCreate, ctx.w)
	} else {
		ctx.data.writeBits(ctx.numBits, ctx.dictionary[ctx.w])
	}
	ctx.decrementEnlargeIn()
}

func compress(uncompressed string) []uint16 {
	ctx := lzCtx{
		dictionary:         map[string]int{},
		dictionaryToCreate: map[string]bool{},
		c:                  rune(0),
		wc:                 "",
		w:                  "",
		enlargeIn:          2, // Compensate for the first entry which should not count
		dictSize:           3,
		numBits:            2,
		result:             []uint16{},
		data:               &lzData{},
	}

	for _, c := range uncompressed {
		ctx.c = c
		if _, ok := ctx.dictionary[string(ctx.c)]; !ok {
			ctx.dictionary[string(ctx.c)] = ctx.dictSize
			ctx.dictSize++
			ctx.dictionaryToCreate[string(ctx.c)] = true
		}

		ctx.wc = ctx.w + string(ctx.c)
		if _, ok := ctx.dictionary[ctx.wc]; ok {
			ctx.w = ctx.wc
		} else {
			ctx.produceW()
			// Add wc to the dictionary.
			ctx.dictionary[ctx.wc] = ctx.dictSize
			ctx.dictSize++
			ctx.w = string(ctx.c)
		}
	}

	// Output the code for w.
	if ctx.w != "" {
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
	res := data.val & uint16(data.position)
	data.position >>= 1
	if data.position == 0 {
		data.position = 32768
		data.val = data.str[data.index]
		data.index++
	}
	//data.val = (data.val << 1);
	if res > 0 {
		return 1
	} else {
		return 0
	}
}

func (data *lzData) readBits(numBits int) uint16 {
	var res = uint16(0)
	var maxpower = int(math.Pow(2, float64(numBits)))
	var power = 1
	for power != maxpower {
		res |= uint16(data.readBit() * power)
		power <<= 1
	}
	return res
}

func decompress(compressed []uint16) string {
	dictionary := map[int]string{}
	enlargeIn := 4
	dictSize := 4
	numBits := 3
	entry := ""
	result := ""
	var w string
	var c uint16
	errorCount := 0
	data := &lzData{
		str:      compressed,
		val:      compressed[0],
		position: 32768,
		index:    1,
	}

	for i := 0; i < 3; i += 1 {
		dictionary[i] = string(rune(i))
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
		return ""
	}
	dictionary[3] = string(rune(c))
	w = string(rune(c))
	result = string(rune(c))

	for {
		c = data.readBits(numBits)

		switch c {
		case 0:
			errorCount++
			if errorCount > 10000 {
				return "Error"
			}
			c = data.readBits(8)
			dictionary[dictSize] = string(rune(c))
			dictSize++
			c = uint16(dictSize - 1)
			enlargeIn--
			break
		case 1:
			c = data.readBits(16)
			dictionary[dictSize] = string(rune(c))
			dictSize++
			c = uint16(dictSize - 1)
			enlargeIn--
			break
		case 2:
			return result
		}

		if enlargeIn == 0 {
			enlargeIn = int(math.Pow(2, float64(numBits)))
			numBits++
		}

		if _, ok := dictionary[int(c)]; ok {
			entry = dictionary[int(c)]
		} else {
			if c == uint16(dictSize) {
				entry = w + string(strings.IndexRune(w, 0))
			} else {
				return "" // This should probably be an error
			}
		}
		result += entry

		// Add w+entry[0] to the dictionary.
		dictionary[dictSize] = w + string(strings.IndexRune(entry, 0))
		dictSize++
		enlargeIn--

		w = entry

		if enlargeIn == 0 {
			enlargeIn = int(math.Pow(2, float64(numBits)))
			numBits++
		}

	}
	return result
}
