package gostikkit

import (
	"encoding/base64"
	"fmt"
	"testing"
)

var key = ``
var lzplain = ``
var plain = ``

func TestCrypt1(t *testing.T) {
	p := []byte(plain)

	p64 := base64.StdEncoding.EncodeToString(p)

	fmt.Println(p64)
}
