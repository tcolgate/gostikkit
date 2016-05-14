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
	DefaultClient.Base, _ = url.Parse("http://paste.scratchbook.ch")
	DefaultClient.hc = http.DefaultClient
}

var rawCall = "/view/raw"

type Client struct {
	Base     *url.URL
	Key      string
	defaults Paste
	hc       *http.Client
}

type Paste struct {
	title     *string
	name      *string
	private   *bool
	lang      *string
	expire    *time.Duration
	replyTo   *string
	decrypted *bytes.Buffer
	key       string
	raw       io.ReadCloser
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

	key := ""
	if u.Fragment != "" {
		key = u.Fragment
	}

	r, err := c.hc.Get(u.String())
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	return &Paste{
		raw: r.Body,
		key: key,
	}, err
}

func Put(p Paste, r io.Reader, encrypt bool) (string, error) {
	return DefaultClient.Put(p, r, encrypt)
}

func (c Client) Put(p Paste, r io.Reader, crypt bool) (string, error) {
	form := url.Values{}

	if p.title != nil {
		form.Add("title", *p.title)
	}

	if p.name != nil {
		form.Add("name", *p.name)
	}

	if p.private != nil && *p.private {
		form.Add("private", "1")
	}

	if p.lang != nil {
		form.Add("lang", *p.lang)
	}

	if p.replyTo != nil {
		form.Add("reply", *p.replyTo)
	}

	if p.expire != nil && (*p.expire).String() != "" {
		form.Add("expire", (*p.expire).String())
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, r)

	key := "sZJf8robYvrQjy5fV3CbDqw7UF5KjVqh"
	salt := []byte(key)
	if crypt {
		lztext := lzjs.CompressToBase64(buf.String())
		ciphertext := encrypt(lztext, key, salt)
		cipherb64 := base64.StdEncoding.EncodeToString(ciphertext)
		buf.Reset()
		buf.Write([]byte(cipherb64))
	}

	form.Add("text", buf.String())

	resp, err := c.hc.PostForm(c.Base.String()+"/api/create", form)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to create paste, " + err.Error())
	}

	buf.Reset()
	io.Copy(buf, resp.Body)

	if crypt {
		url := buf.String()
		url = fmt.Sprintf("%s#%s", strings.Replace(url, "\n", "", -1), key)
		buf.Write([]byte(url))
	}
	return buf.String(), nil
}

var DefaultClient = &Client{}

func NewClient() *Client {
	c := *DefaultClient
	return &c
}

type option func(c *Client) option
type pasteoption func(c *Paste) pasteoption

func Option(opts ...option) (previous option) {
	return DefaultClient.Option(opts...)
}

func (c *Client) Option(opts ...option) (previous option) {
	for _, opt := range opts {
		previous = opt(c)
	}
	return previous
}

func DefaultExpire(t time.Duration) option {
	return func(c *Client) option {
		previous := c.defaults.expire
		c.defaults.expire = &t
		return DefaultExpire(*previous)
	}
}

func DefaultName(n string) option {
	return func(c *Client) option {
		previous := c.defaults.name
		c.defaults.name = &n
		return DefaultName(*previous)
	}
}

func DefaultPrivate(p bool) option {
	return func(c *Client) option {
		previous := c.defaults.private
		c.defaults.private = &p
		return DefaultPrivate(*previous)
	}
}

func (c *Paste) Option(opts ...pasteoption) (previous pasteoption) {
	for _, opt := range opts {
		previous = opt(c)
	}
	return previous
}

func Expire(t time.Duration) pasteoption {
	return func(p *Paste) pasteoption {
		previous := p.expire
		p.expire = &t
		return Expire(*previous)
	}
}

func Name(n string) pasteoption {
	return func(p *Paste) pasteoption {
		previous := p.name
		p.name = &n
		return Name(*previous)
	}
}

func Private(b bool) pasteoption {
	return func(p *Paste) pasteoption {
		previous := p.private
		p.private = &b
		return Private(*previous)
	}
}

func (p *Paste) Read(bs []byte) (int, error) {
	if p.key != "" {
		if p.decrypted == nil {
			//Don't currently support lzjs using a reader
			buf := &bytes.Buffer{}
			io.Copy(buf, p.raw)
			clean := strings.Replace(string(buf.Bytes()), "\n", "", -1)

			ciphertext, err := base64.StdEncoding.DecodeString(clean)
			if err != nil {
				return 0, err
			}
			lztext := decrypt(ciphertext, p.key)
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
