// Package assetfs provides an http.FileSystem for go-bindata assets.
package assetfs

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"sort"
	"time"
)

// AssetStore provides access to assets generated by go-bindata.
type AssetStore struct {
	Names func() []string
	Data  func(filename string) ([]byte, error)
	Info  func(filename string) (os.FileInfo, error)
}

// AssetFS is a static file server that implements the http.FileSystem interface.
type AssetFS struct {
	m map[string]*assetData
}

type assetData struct {
	b    []byte
	info os.FileInfo
}

// New creates an AssetFS from a list of files and their bytes.
// Pattern is a glob from path.Match, and strip true means to strip
// off the pattern (e.g. dir prefix) from the file name.
func New(f *AssetStore) (*AssetFS, error) {
	m := make(map[string]*assetData)
	for _, fn := range f.Names() {
		info, err := f.Info(fn)
		if err != nil {
			return nil, err
		}
		data, err := f.Data(fn)
		if err != nil {
			return nil, err
		}
		dir, file := path.Split(fn)
		if file == "index.html" {
			dir = AddPrefix("/", path.Clean(dir))
			m[dir] = &assetData{b: []byte{}, info: dirInfo(dir)}
		}
		k := AddPrefix("/", fn)
		m[k] = &assetData{b: data, info: info}
	}
	return &AssetFS{m}, nil
}

// Open implements the http.FileSystem interface.
func (fs *AssetFS) Open(name string) (http.File, error) {
	f, exists := fs.m[name]
	if !exists {
		return nil, os.ErrNotExist
	}
	return newFile(f.b, f.info), nil
}

// Files returns a list of files in the AssetFS.
func (fs *AssetFS) Files() []string {
	f := make([]string, len(fs.m))
	i := 0
	for fn := range fs.m {
		f[i] = fn
	}
	sort.Strings(f)
	return f
}

// Len returns the length of the given file in the AssetFS.
func (fs *AssetFS) Len(name string) int {
	f, exists := fs.m[name]
	if !exists {
		return 0
	}
	return len(f.b)
}

func newFile(data []byte, info os.FileInfo) *assetFile {
	r := bytes.NewReader(data)
	return &assetFile{
		Reader:   r,
		FileInfo: info,
	}
}

type assetFile struct {
	*bytes.Reader
	os.FileInfo
}

func (f *assetFile) Close() error {
	return nil
}

func (f *assetFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, os.ErrInvalid
}

func (f *assetFile) Stat() (os.FileInfo, error) {
	return f.FileInfo, nil
}

type dirInfo string

func (d dirInfo) Name() string       { return string(d) }
func (d dirInfo) Size() int64        { return 0 }
func (d dirInfo) Mode() os.FileMode  { return 0 }
func (d dirInfo) ModTime() time.Time { return time.Now() }
func (d dirInfo) IsDir() bool        { return true }
func (d dirInfo) Sys() interface{}   { return nil }
