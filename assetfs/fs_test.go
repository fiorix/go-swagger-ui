//go:generate go-bindata -pkg=internal -o internal/files.go fs.go

package assetfs

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/fiorix/go-swagger-ui/assetfs/internal"
)

func TestAssetFS(t *testing.T) {
	assets := &AssetStore{
		Names: internal.AssetNames,
		Data:  internal.Asset,
		Info:  internal.AssetInfo,
	}
	fs, err := New(assets)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(fs))
	s := httptest.NewServer(mux)
	defer s.Close()
	for _, fn := range fs.Files() {
		l := fs.Len(fn)
		resp, err := http.Get(s.URL + "/" + filepath.Base(fn))
		if err != nil {
			t.Errorf("Test %q: %v", fn, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Test %q: status %q", fn, resp.Status)
		}
		if resp.ContentLength != int64(l) {
			t.Errorf("Test %q:\nWant: %d\nHave: %d",
				fn, l, resp.ContentLength)
		}
	}
}
