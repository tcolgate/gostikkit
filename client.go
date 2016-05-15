// Copyright (c) 2016 Tristan Colgate-McFarlane
//
// This file is part of gostikkit.
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

package gostikkit

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tcolgate/gostikkit/lzjs"
)

func init() {
	DefaultClient.hc = http.DefaultClient
}

var rawCall = "/view/raw"

type Client struct {
	Base *url.URL
	Key  string

	Paste //The default post to use

	hc *http.Client
}

const (
	ExpireNever        = time.Duration(-1)
	ExpireAfterReading = time.Duration(-2)
)

type Paste struct {
	Title   string
	Author  string
	Private bool
	Lang    string
	Expire  time.Duration
	ReplyTo string

	decrypted *bytes.Buffer
	ckey      string
	raw       io.ReadCloser
}

func (c Client) New() Paste {
	return c.Paste
}

func Get(id string) (*Paste, error) {
	return DefaultClient.Get(id)
}

func (c Client) Get(id string) (*Paste, error) {
	u, err := url.Parse(id)
	if err != nil || u.Host == "" {
		ustr := fmt.Sprintf("%v/%v/%v", c.Base, rawCall, id)
		u, err = url.Parse(ustr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "")
			os.Exit(1)
		}
	}

	ups := strings.Split(u.Path, "/")
	if len(ups) >= 2 {
		call := ups[len(ups)-2 : len(ups)-1][0]
		if call != "raw" {
			if call == "view" {
				prefix := ups[:len(ups)-1]
				id := ups[len(ups)-1]
				newparts := append(prefix, "raw", id)
				u.Path = strings.Join(newparts, "/")
			}
		}
	}

	ckey := ""
	if u.Fragment != "" {
		ckey = u.Fragment
	}

	if u.Query().Get("apikey") == "" && c.Key != "" {
		vs := u.Query()
		vs.Add("apikey", c.Key)
		u.RawQuery = vs.Encode()
	}
	r, err := c.hc.Get(u.String())
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	return &Paste{
		raw:  r.Body,
		ckey: ckey,
	}, err
}

func Put(p Paste, r io.Reader, encrypt bool) (string, error) {
	return DefaultClient.Put(p, r, encrypt)
}

func (c Client) Put(p Paste, r io.Reader, crypt bool) (string, error) {
	form := url.Values{}

	if p.Title != "" {
		form.Add("title", p.Title)
	}

	if p.Author != "" {
		form.Add("name", p.Author)
	}

	if p.Private {
		form.Add("private", "1")
	}

	if p.Lang != "" {
		form.Add("lang", p.Lang)
	}

	if p.ReplyTo != "" {
		form.Add("reply", p.ReplyTo)
	}

	if p.Expire != time.Duration(0) {
		if p.Expire == ExpireNever {
			form.Add("expire", "0")
		} else if p.Expire == ExpireAfterReading {
			form.Add("expire", "burn")
		} else {
			form.Add("expire", fmt.Sprintf("%d", p.Expire/time.Minute))
		}
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, r)

	key := genChars(32)
	if crypt {
		lztext := lzjs.CompressToBase64(buf.String())
		ciphertext := encrypt(lztext, string(key))
		cipherb64 := base64.StdEncoding.EncodeToString(ciphertext)
		buf.Reset()
		buf.Write([]byte(cipherb64))
	}

	form.Add("text", buf.String())

	rurl, err := url.Parse(c.Base.String() + "/api/create")
	if err != nil {
		return "", errors.New("failed consutructing reqyest url" + err.Error())
	}

	if rurl.Query().Get("apikey") == "" && c.Key != "" {
		vs := rurl.Query()
		vs.Add("apikey", c.Key)
		rurl.RawQuery = vs.Encode()
	}
	resp, err := c.hc.PostForm(rurl.String(), form)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to create paste, " + err.Error())
	}

	buf.Reset()
	io.Copy(buf, resp.Body)
	url := buf.String()
	url = strings.Replace(url, "\n", "", -1)

	if crypt {
		url = fmt.Sprintf("%s#%s", url, key)
	}
	return url, nil
}

var DefaultClient = &Client{}

func NewClient() *Client {
	c := *DefaultClient
	return &c
}

func (p *Paste) Read(bs []byte) (int, error) {
	if p.ckey != "" {
		if p.decrypted == nil {
			//Don't currently support lzjs using a reader
			buf := &bytes.Buffer{}
			io.Copy(buf, p.raw)
			clean := strings.Replace(string(buf.Bytes()), "\n", "", -1)

			ciphertext, err := base64.StdEncoding.DecodeString(clean)
			if err != nil {
				return 0, err
			}
			lztext := decrypt(ciphertext, p.ckey)
			plain, err := lzjs.DecompressFromBase64(string(lztext))
			if err != nil {
				return 0, err
			}

			p.decrypted = bytes.NewBuffer([]byte(plain))
		}
		return p.decrypted.Read(bs)
	}
	return p.raw.Read(bs)
}

func (p *Paste) Close() error {
	return p.raw.Close()
}
