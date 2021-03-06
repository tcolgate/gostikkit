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

package lzjs

import (
	"reflect"
	"testing"
	"testing/quick"
	"unicode/utf16"
)

var testingValues map[string][]uint16 = map[string][]uint16{
	"A":                          {8336},
	"AA":                         {8370},
	"AB":                         {8322, 4608},
	"ABC":                        {8322, 4290, 16384},
	"ABCABC":                     {8322, 4290, 42560},
	"ABCACABCACABC":              {8322, 4290, 49641, 42560},
	"A☺A☺":                       {8354, 58149, 16384},
	"ABCλλBCλAC":                 {8322, 4290, 36316, 481, 50752},
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ": {8322, 4290, 544, 41478, 8418, 144, 9344, 41997, 8217, 712, 3648, 61952, 10248, 40978, 32970, 168, 2720, 6784, 59904, 26633, 40982, 36864},
}

var gatsby = `In my younger and more vulnerable years my father gave me some advice that I’ve been turning over in my mind ever since.
“Whenever you feel like criticizing any one,” he told me, “just remember that all the people in this world haven’t had the advantages that you’ve had.”
He didn’t say any more, but we’ve always been unusually communicative in a reserved way, and I understood that he meant a great deal more than that. In consequence, I’m inclined to reserve all judgments, a habit that has opened up many curious natures to me and also made me the victim of not a few veteran bores. The abnormal mind is quick to detect and attach itself to this quality when it appears in a normal person, and so it came about that in college I was unjustly accused of being a politician, because I was privy to the secret griefs of wild, unknown men. Most of the confidences were unsought — frequently I have feigned sleep, preoccupation, or a hostile levity when I realized by some unmistakable sign that an intimate revelation was quivering on the horizon; for the intimate revelations of young men, or at least the terms in which they express them, are usually plagiaristic and marred by obvious suppressions. Reserving judgments is a matter of infinite hope. I am still a little afraid of missing something if I forget that, as my father snobbishly suggested, and I snobbishly repeat, a sense of the fundamental decencies is parcelled out unequally at birth.
And, after boasting this way of my tolerance, I come to the admission that it has a limit. Conduct may be founded on the hard rock or the wet marshes, but after a certain point I don’t care what it’s founded on. When I came back from the East last autumn I felt that I wanted the world to be in uniform and at a sort of moral attention forever; I wanted no more riotous excursions with privileged glimpses into the human heart. Only Gatsby, the man who gives his name to this book, was exempt from my reaction — Gatsby, who represented everything for which I have an unaffected scorn. If personality is an unbroken series of successful gestures, then there was something gorgeous about him, some heightened sensitivity to the promises of life, as if he were related to one of those intricate machines that register earthquakes ten thousand miles away. This responsiveness had nothing to do with that flabby impressionability which is dignified under the name of the “creative temperament.”— it was an extraordinary gift for hope, a romantic readiness such as I have never found in any other person and which it is not likely I shall ever find again. No — Gatsby turned out all right at the end; it is what preyed on Gatsby, what foul dust floated in the wake of his dreams that temporarily closed out my interest in the abortive sorrows and short-winded elations of men.`

func TestCompDecomp(t *testing.T) {
	for k := range testingValues {
		compd := Compress(k)
		result, err := Decompress(compd)
		if err != nil {
			t.Errorf("decompress threw error, %v", err)
		}
		if !reflect.DeepEqual(result, k) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", result, k)
		}
	}
}

func TestDecompComp(t *testing.T) {
	for _, v := range testingValues {
		decom, err := Decompress(v)
		if err != nil {
			t.Errorf("decompress threw error, %v", err)
		}
		result := Compress(decom)
		if !reflect.DeepEqual(result, v) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", result, v)
		}
	}
}

func TestCompDecompGatsby(t *testing.T) {
	s := gatsby
	compd := Compress(s)
	result, err := Decompress(compd)
	if err != nil {
		t.Errorf("decompress threw error, %v", err)
	}
	if result != s {
		t.Errorf("Result should be :\n  %v\n instead of :\n  %v", result, s)
	}
}

func TestQuickCompDecomp(t *testing.T) {
	f := func(cs []uint16) bool {
		x := utf16.Decode(cs)
		y, err := Decompress(Compress(string(x)))
		if err != nil {
			t.Fatal("err ", err)
		}

		return err == nil && string(x) == y
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestCompress(t *testing.T) {
	for k, v := range testingValues {
		result := Compress(k)
		if !reflect.DeepEqual(result, v) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", v, result)
		}
	}
}

func TestDecompress(t *testing.T) {
	for k, v := range testingValues {
		result, err := Decompress(v)
		if err != nil {
			t.Errorf("decompress threw error, %v", err)
		}
		if !reflect.DeepEqual(result, k) {
			t.Errorf("Result should be :\n  %v\n instead of :\n  %v", k, result)
		}
	}
}

func BenchmarkCompress(b *testing.B) {
	s := gatsby
	bs := len(s)
	b.SetBytes(int64(bs))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Compress(s)
	}
}

func BenchmarkDecompress(b *testing.B) {
	s := gatsby
	d := Compress(s)
	bs := len(d) * 2
	b.SetBytes(int64(bs))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Decompress(d)
	}
}

func TestCompressToBase64(t *testing.T) {
	in := "hello"
	res := CompressToBase64(in)

	res2, err := DecompressFromBase64(res)
	if err != nil {
		t.Fatal("error ", err)
	}

	if res2 != in {
		t.Fatal("input and output did not match")
	}
}
