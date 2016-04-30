package lzjs

import (
	"reflect"
	"testing"
)

var testingValues map[string][]uint16 = map[string][]uint16{
	"A":             {},
	"AB":            {},
	"ABC":           {},
	"ABCABC":        {},
	"ABCACABCACABC": {},
}

func TestDecompComp(t *testing.T) {
	for k := range testingValues {
		result := decompress(compress(k))
		if !reflect.DeepEqual(result, k) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", k, result)
		}
	}
}

func TestCompDecomp(t *testing.T) {
	for _, v := range testingValues {
		result := compress(decompress(v))
		if !reflect.DeepEqual(result, v) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", v, result)
		}
	}
}

// func TestDecompress(t *testing.T) {
// for k, v := range testingValues {
// 	result := decompress(v)
// 	/*
// 		result, err := decompress(v)
// 		if err != nil {
// 			t.Errorf("Unexpected error", err)
// 		}
// 	*/
// 	if result != k {
// 		t.Errorf("Result should be :\n", v, "\n instead of :\n", result)
// 	}
// }
//}

func TestCecompress(t *testing.T) {
	for k, v := range testingValues {
		result := compress(k)
		/*
			result, err := decompress(v)
			if err != nil {
				t.Errorf("Unexpected error", err)
			}
		*/
		if !reflect.DeepEqual(result, v) {
			t.Errorf("Result should be :\n", v, "\n instead of :\n", result)
		}
	}
}

/*
func BenchmarkDecompress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for k := range testingValues {
			DecompressFromEncodedUriComponent(k)
		}
	}
}

func TestCompress(t *testing.T) {
	for k, v := range testingValues {
		result, err := CompressToBase64(v)
		if err != nil {
			t.Errorf("Unexpected error", err)
		}
		if result != k {
			t.Errorf("Result should be :\n", v, "\n instead of :\n", result)
		}
	}
}

func BenchmarkCompress(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range testingValues {
			CompressToBase64(v)
		}
	}
}
*/
