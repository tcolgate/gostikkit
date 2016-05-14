package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/tcolgate/gostikkit"
)

var (
	urlStr = flag.String("url", "http://localhost:80", "Post contents of this file")

	author  = flag.String("author", "", "Author of the post")
	name    = flag.String("name", "", "Default name  of the post")
	expire  = flag.String("expire", "", "Expiration time, in minutes (or \"never\", or \"burn\")")
	encrypt = flag.Bool("encrypt", false, "Encrypt the post")
	file    = flag.String("file", "", "Post contents of this file")
)

func main() {
	flag.Parse()

	if *urlStr != "" {
		var err error
		gostikkit.DefaultClient.Base, err = url.Parse(*urlStr)
		if err != nil {
			fmt.Println("could not parse URL %v, %v", *urlStr, err.Error())
		}
	}

	if len(flag.Args()) == 0 || *file != "" {
		var r = os.Stdin
		var err error
		if *urlStr == "" {
			fmt.Printf("A url must be specific for pasting\n")
			os.Exit(1)
		}
		if *file != "-" && *file != "" {
			r, err = os.Open(*file)
			if err != nil {
				fmt.Printf("could not open %v, %v", *file, err.Error())
			}
			os.Exit(1)
		}

		p := gostikkit.Paste{}
		purl, err := gostikkit.Put(p, r, *encrypt)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		fmt.Println(purl)
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
