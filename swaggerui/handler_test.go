package swaggerui

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fiorix/go-swagger-ui/assetfs"
)

func TestHandler(t *testing.T) {
	cases := []struct {
		Request *http.Request
		Code    int
	}{
		{
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "/"},
			},
			Code: http.StatusSeeOther,
		},
		{
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: SpecFile},
			},
			Code: http.StatusOK,
		},
		{
			Request: &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: BasePath},
			},
			Code: http.StatusOK,
		},
	}
	for _, prefix := range []string{"/", "/foobar/"} {
		f := Handler(prefix, strings.NewReader("{}"))
		for _, tc := range cases {
			w := &httptest.ResponseRecorder{}
			tc.Request.URL.Path = assetfs.AddPrefix(prefix, tc.Request.URL.Path)
			f.ServeHTTP(w, tc.Request)
			if w.Code != tc.Code {
				t.Logf("headers: %#v", w.Header())
				t.Fatalf("unexpected code for %s: want %d, have %d",
					tc.Request.URL.Path, tc.Code, w.Code)
			}
		}
	}
}
