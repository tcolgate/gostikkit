package lzjs

import (
	"reflect"
	"testing"
)

var testingValues map[string][]uint16 = map[string][]uint16{
	"A":             {8336},
	"AB":            {8322, 4608},
	"ABC":           {8322, 4290, 16384},
	"ABCABC":        {8322, 4290, 42560},
	"ABCACABCACABC": {8322, 4290, 49641, 42560},
	"A☺A☺":          {8354, 58149, 16384},
	"ABCλλBCλAC":    {8322, 4290, 36316, 481, 50752},
}

func TestDecompComp(t *testing.T) {
	for k := range testingValues {
		compd := compress(k)
		result := decompress(compd)
		if !reflect.DeepEqual(result, k) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", result, k)
		}
	}
}

func TestCompDecomp(t *testing.T) {
	for _, v := range testingValues {
		result := compress(decompress(v))
		if !reflect.DeepEqual(result, v) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", result, v)
		}
	}
}

/*
func TestQuickCompDecomp(t *testing.T) {
	f := func(x string) bool {
		y := decompress(compress(x))
		return x == y
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
*/

func TestCompress(t *testing.T) {
	for k, v := range testingValues {
		result := compress(k)
		if !reflect.DeepEqual(result, v) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", result, v)
		}
	}
}

func TestDecompress(t *testing.T) {
	for k, v := range testingValues {
		result := decompress(v)
		if !reflect.DeepEqual(result, k) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", result, k)
		}
	}
}

func BenchmarkCompress(b *testing.B) {
	//sr, _ := quick.Value(reflect.TypeOf("string"), nil)
	//s := sr.String()
	s := "hello"
	bs := len(s)
	b.SetBytes(int64(bs))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		compress(s)
	}
}

func BenchmarkDecompress(b *testing.B) {
	// run the Fib function b.N times
	// sr, _ := quick.Value(reflect.TypeOf("string"), nil)
	// s := sr.String()
	s := "hello"
	d := compress(s)
	bs := len(d) * 2
	b.SetBytes(int64(bs))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		decompress(d)
	}
}
