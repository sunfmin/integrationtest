package integrationtest

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

var Verbose bool

type cookieJar struct {
	cookies []*http.Cookie
}

type multipartBuilder func(w *multipart.Writer)

func (tc *cookieJar) find(cookie *http.Cookie) (at int, r *http.Cookie) {
	for i, c := range tc.cookies {
		if c.Name == cookie.Name {
			at = i
			r = c
			return
		}
	}
	return
}

func (tc *cookieJar) String() (r string) {
	for _, c := range tc.cookies {
		r = r + c.Name + " => " + c.Value + ", "
	}
	return
}

func (tc *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, c := range cookies {
		at, fc := tc.find(c)
		if fc != nil {
			tc.cookies[at] = c
		} else {
			tc.cookies = append(tc.cookies, c)
		}
	}
	if Verbose {
		log.Printf("Set cookie: %s\n", tc)
	}
	return
}

func (tc *cookieJar) Cookies(u *url.URL) []*http.Cookie {
	if Verbose {
		log.Printf("Cookies: %s\n", tc)
	}
	return tc.cookies
}

type Session struct {
	Client *http.Client
}

func NewSession() (s *Session) {
	s = &Session{}
	s.Client = &http.Client{}
	s.Client.Jar = &cookieJar{}
	return
}

func (s *Session) Get(u string) (r *http.Response, err error) {
	if Verbose {
		log.Printf("Get %s\n", u)
	}
	r, err = s.Client.Get(u)

	// if Verbose {
	// 	log.Printf("After Get, Cookie in Response: %+v", r.Header)
	// }

	// if c := r.Cookies(); len(c) > 0 {
	// 	s.Client.Jar.SetCookies(nil, c)
	// }

	return
}

func (s *Session) Post(u string, bodyType string, body io.Reader) (r *http.Response, err error) {
	if Verbose {
		log.Printf("Post %s\n", u)
	}

	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", bodyType)
	c := s.Client

	for _, cookie := range c.Jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	r, err = c.Do(req)
	if err == nil && c.Jar != nil {
		c.Jar.SetCookies(req.URL, r.Cookies())
	}
	return r, err
}

func (s *Session) PostMultipart(u string, mb multipartBuilder) (r *http.Response, err error) {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	mb(mw)
	mw.Close()
	return s.Post(u, fmt.Sprintf("multipart/form-data; boundary=%s", mw.Boundary()), &b)
}

func (s *Session) PostForm(u string, data url.Values) (r *http.Response, err error) {
	return s.Post(u, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (s *Session) Head(u string) (r *http.Response, err error) {
	if Verbose {
		log.Printf("Head %s\n", u)
	}
	return s.Client.Head(u)
}

func Must(in *http.Response, err error) (r *http.Response) {
	if err != nil {
		panic(err)
	}
	r = in
	return
}
