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
}

func TestDecompComp(t *testing.T) {
	for k := range testingValues {
		compd := compress(k)
		result := decompress(compd)
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

func TestCompress(t *testing.T) {
	for k, v := range testingValues {
		result := compress(k)
		if !reflect.DeepEqual(result, v) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", k, result)
		}
	}
}

func TestDecompress(t *testing.T) {
	for k, v := range testingValues {
		result := decompress(v)
		if !reflect.DeepEqual(result, k) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", v, result)
		}
	}
}
