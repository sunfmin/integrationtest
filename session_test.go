package integrationtest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func theMux() (sm *http.ServeMux) {
	sm = http.NewServeMux()
	sm.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Started %s \"%s\" for %s\n", r.Method, r.RequestURI, r.RemoteAddr)
		cookie := &http.Cookie{
			Name:  "name",
			Value: "felix",
		}
		http.SetCookie(w, cookie)
	})

	sm.HandleFunc("/updateprofile", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Started %s \"%s\" for %s\n", r.Method, r.RequestURI, r.RemoteAddr)
		cookie := &http.Cookie{
			Name:  "age",
			Value: "23",
		}
		http.SetCookie(w, cookie)
	})

	sm.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Started %s \"%s\" for %s\n", r.Method, r.RequestURI, r.RemoteAddr)
		name, _ := r.Cookie("name")
		fmt.Fprintf(w, "%s", name)
	})
	return
}

func TestHoldingCookies(t *testing.T) {
	Verbose = true
	ts := httptest.NewServer(theMux())
	defer ts.Close()

	s := NewSession()

	fmt.Println("before login: ", s.Client.Jar.Cookies(nil))

	loginRes, _ := s.Get(ts.URL + "/login")
	fmt.Println("after login: ", loginRes.Cookies())

	upres, _ := s.Get(ts.URL + "/updateprofile")
	fmt.Println("after updateprofile: ", upres.Cookies())

	res, _ := s.PostForm(ts.URL+"/account", url.Values{})
	fmt.Println("response cookies: ", res.Cookies())

	fmt.Println("after account: ", s.Client.Jar.Cookies(nil))

	b, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))

	if strings.Index(string(b), "felix") < 0 {
		t.Errorf("response body should contain cookie value: %+v", string(b))
	}
}
