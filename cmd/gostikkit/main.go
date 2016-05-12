package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tcolgate/gostikkit"
)

var (
	author = flag.String("author", "", "Author of the post")
	name   = flag.String("name", "", "Default name  of the post")
	expire = flag.String("expire", "", "Expiration time, in minutes (or \"never\", or \"burn\")")
	file   = flag.String("file", "", "Post contents of this file")
	url    = flag.String("url", "", "Post contents of this file")
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 || *file != "" {
		var r = os.Stdin
		var err error
		if *url == "" {
			fmt.Printf("A url must be specific for pasting\n")
			os.Exit(1)
		}
		if *file != "-" && *file != "" {
			r, err = os.Open(*file)
			if err != nil {
				fmt.Println("could not open %v, %v", *file, err.Error())
				os.Exit(1)
			}
		}
		io.Copy(os.Stdout, r)
		return
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
