package integrationtest

import (
	"io"
	"net/http"
	"net/url"
)

type CookieJar struct {
	cookies []*http.Cookie
}

func (tc *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	tc.cookies = cookies
	return
}

func (tc *CookieJar) Cookies(u *url.URL) []*http.Cookie {
	return tc.cookies
}

type Session struct {
	Jar *CookieJar
}

func NewSession() (s *Session) {
	return &Session{&CookieJar{}}
}
func (s *Session) ClientWithJar() (client *http.Client) {
	client = &http.Client{}
	client.Jar = s.Jar
	return
}

func (s *Session) Get(url string) (r *http.Response, err error) {
	return s.ClientWithJar().Get(url)
}

func (s *Session) Post(url string, bodyType string, body io.Reader) (r *http.Response, err error) {
	return s.ClientWithJar().Post(url, bodyType, body)
}

func (s *Session) PostForm(url string, data url.Values) (r *http.Response, err error) {
	return s.ClientWithJar().PostForm(url, data)
}

func (s *Session) Head(url string) (r *http.Response, err error) {
	return s.ClientWithJar().Head(url)
}
