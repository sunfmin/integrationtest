package integrationtest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testMux() (sm *http.ServeMux) {
	sm = http.NewServeMux()
	sm.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:  "name",
			Value: "felix",
		}
		http.SetCookie(w, cookie)
	})
	sm.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		name, _ := r.Cookie("name")
		fmt.Fprintf(w, "%s", name.String())
	})
	return
}

func TestHoldingCookies(t *testing.T) {
	ts := httptest.NewServer(testMux())
	defer ts.Close()

	s := NewSession()

	s.Get(ts.URL + "/login")

	res, _ := s.Get(ts.URL + "/account")
	b, _ := ioutil.ReadAll(res.Body)

	if strings.Index(string(b), "felix") < 0 {
		t.Errorf("response body should contain cookie value: %+v", string(b))
	}
}
