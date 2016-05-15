package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/tcolgate/gostikkit"
)

var (
	urlStr = flag.String("url", "http://localhost:80", "Post contents of this file")

	author  = flag.String("author", "", "Author of the post")
	title   = flag.String("title", "", "Title of the post")
	lang    = flag.String("lang", "", "The language to render the post as")
	expire  = flag.String("expire", "", "Expiration time, in minutes (or \"never\", or \"burn\")")
	encrypt = flag.Bool("encrypt", false, "Encrypt the post")
	file    = flag.String("file", "", "Post contents of this file")
)

func main() {
	var err error
	flag.Parse()

	if *urlStr != "" {
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
				os.Exit(1)
			}
		}

		p := gostikkit.Paste{}

		if *title != "" {
			p.Title = *title
		} else if *file != "" {
			p.Title = path.Base(*file)
		}

		if *author != "" {
			p.Author = *author
		} else {
			u, err := user.Current()
			if err == nil {
				p.Author = u.Username
			}
			h, err := os.Hostname()
			if p.Author != "" && err == nil {
				p.Author = p.Author + "@" + h
			}
		}

		if *lang != "" {
			p.Lang = *lang
		} else if *file != "" {
			fn := path.Base(*file)
			fps := strings.Split(fn, ".")
			if len(fps) > 1 {
				p.Lang = fps[len(fps)-1]
			}
		}

		if *expire != "" {
			if *expire == "burn" {
				p.Expire = gostikkit.ExpireAfterReading
			} else if *expire == "never" || *expire == "0" {
				p.Expire = gostikkit.ExpireNever
			} else {
				ms, err := strconv.Atoi(*expire)
				if err != nil {
					fmt.Printf("expiration must be a number of minute, \"never\", or \"burn\"")
					os.Exit(1)
				}
				p.Expire = time.Minute * time.Duration(ms)
			}
		}

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
