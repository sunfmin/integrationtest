package integrationtest

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var Verbose bool

type cookieJar struct {
	cookies []*http.Cookie
}

func (tc *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	tc.cookies = cookies
	return
}

func (tc *cookieJar) Cookies(u *url.URL) []*http.Cookie {
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
	return s.Client.Post(u, bodyType, body)
}

func (s *Session) PostForm(u string, data url.Values) (r *http.Response, err error) {
	if Verbose {
		log.Printf("PostForm %s\n", u)
	}
	return s.Client.PostForm(u, data)
}

func (s *Session) Head(u string) (r *http.Response, err error) {
	if Verbose {
		log.Printf("Head %s\n", u)
	}
	return s.Client.Head(u)
}
