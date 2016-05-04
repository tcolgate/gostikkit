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
	"io"
	"net/http"
	"net/url"
	"time"
)

func init() {
	DefaultClient.hc = http.DefaultClient
}

type Client struct {
	Base     url.URL
	Key      string
	defaults Paste
	hc       *http.Client
}

type Paste struct {
	title   *string
	name    *string
	private *bool
	lang    *string
	expire  *time.Duration
	replyTo *string
	io.ReadCloser
}

func Get(id string) (Paste, error) {
	return DefaultClient.Get(id)
}

func (c Client) Get(id string) (Paste, error) {
	r, err := c.hc.Get(id)
	return Paste{ReadCloser: r.Body}, err
}

func Put(p Paste, encrypt bool) (string, error) {
	return DefaultClient.Put(p, encrypt)
}

func (c Client) Put(p Paste, encrypt bool) (string, error) {
	//	form := url.Values{}
	// form.Add("title", p.Title)
	// form.Add("name", p.Name)
	// //form.Add("private", p.Private)
	// form.Add("lang", p.Lang)
	// form.Add("expire", p.Expire.String())
	// form.Add("reply", p.ReplyTo)
	return "", nil
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
