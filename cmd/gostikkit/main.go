package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/tcolgate/gostikkit"
)

func main() {
	flag.Parse()

	if len(os.Args[1:]) == 0 {
		log.Println("no args, upload")
	}

	for _, a := range os.Args[1:] {
		p, err := gostikkit.Get(a)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = io.Copy(os.Stdout, p)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
