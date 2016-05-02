//go:generate sh -c "go-bindata -pkg=internal -o internal/files.go `find third_party -type d`"

package swaggerui

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fiorix/go-swagger-ui/assetfs"
	"github.com/fiorix/go-swagger-ui/swaggerui/internal"
)

// BasePath is the base path of the embedded swagger-ui.
const BasePath = "/third_party/swagger-ui/"

// SpecFile is the name of the swagger JSON file to serve.
const SpecFile = "/swagger.json"

// Handler returns an HTTP handler that serves the
// swagger-ui under /third_party/swagger-ui.
func Handler(basepath string, data io.ReadSeeker) http.Handler {
	if basepath == "" {
		basepath = "/"
	}
	as := &assetfs.AssetStore{
		Names: internal.AssetNames,
		Data:  internal.Asset,
		Info:  internal.AssetInfo,
	}
	fs, err := assetfs.New(as)
	if err != nil {
		panic(fmt.Sprintf("failed to create static fs: %v", err))
	}
	mux := http.NewServeMux()
	fsh := http.FileServer(http.FileSystem(fs))
	if basepath != "/" {
		fsh = http.StripPrefix(basepath, fsh)
	}
	p := assetfs.AddPrefix(basepath, BasePath)
	f := assetfs.AddPrefix(basepath, SpecFile)
	mux.HandleFunc(basepath, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == basepath {
			http.Redirect(w, r, p+"?url="+f, http.StatusSeeOther)
			return
		}
		fsh.ServeHTTP(w, r)
	})
	mux.Handle(f, &handler{modTime: time.Now(), body: data})
	return mux
}

type handler struct {
	modTime time.Time
	body    io.ReadSeeker
}

func (f *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, SpecFile, f.modTime, f.body)
}
