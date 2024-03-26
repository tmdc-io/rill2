package web

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
)

//go:embed all:embed
var distFS embed.FS

// Handler serves an web-local UI.
func StaticHandler() http.Handler {
	uiAssetFS := newUIAssetFS()
	return gziphandler.GzipHandler(http.FileServer(uiAssetFS))
}

// Check if web-local dist static UI is exists, If not server the default index.html page.
func newUIAssetFS() http.FileSystem {
	_, err := distFS.ReadFile("embed/dist/index.html")
	if os.IsNotExist(err) {
		return assetFS(distFS, "embed")
	}
	return assetFS(distFS, "embed/dist")
}

// Get the subtree of the embedded files with `embed` directory as a root.
func assetFS(embeddedFS embed.FS, dir string) http.FileSystem {
	subFS, err := fs.Sub(embeddedFS, dir)
	if err != nil {
		panic(fmt.Errorf("fs embed: %w", err))
	}

	return &SPARoutingFS{FileSystem: http.FS(subFS)}
}

type SPARoutingFS struct {
	FileSystem http.FileSystem
}

//func (spaFS *SPARoutingFS) Open(name string) (http.File, error) {
//	file, err := spaFS.FileSystem.Open(name)
//	if err == nil {
//		return file, nil
//	}
//
//	if errors.Is(err, fs.ErrNotExist) {
//		file, err = spaFS.FileSystem.Open("index.html")
//		return file, err
//	}
//
//	return nil, err
//}

func (spaFS *SPARoutingFS) Open(name string) (http.File, error) {
	ext := filepath.Ext(name)
	var err error
	if ext == ".css" || ext == ".js" || ext == ".html" {
		data, err := spaFS.FileSystem.Open(name)
		if err != nil {
			return nil, err
		}
		defer data.Close()

		fileBytes, err := ioutil.ReadAll(data)
		if err != nil {
			return nil, err
		}

		// Replace "base-path" with the desired value in memory
		newData := strings.ReplaceAll(string(fileBytes), "dataos-basepath", strings.Trim(os.Getenv("BASE_PATH"), "/"))

		fileInfo := &mockFileInfo{
			size:    int64(len(newData)), // Approximate size based on modified data
			modTime: time.Now(),          // Use current time for modification
			isDir:   false,               // Assume it's not a directory
			name:    name,                // Use the original filename
		}

		customFile := &CustomFile{
			Reader:   bytes.NewReader([]byte(newData)),
			FileInfo: fileInfo, // Get FileInfo from the original file
		}

		return customFile, nil

	} else {
		file, err := spaFS.FileSystem.Open(name)
		if err == nil {
			return file, nil
		}
	}
	if errors.Is(err, fs.ErrNotExist) {
		file, err := spaFS.FileSystem.Open("index.html")
		return file, err
	}
	return nil, err
}

// CustomFile implements the http.File interface for in-memory modifications
type CustomFile struct {
	*bytes.Reader
	FileInfo os.FileInfo
}

// Seek implements the Seek method of the http.File interface
func (f *CustomFile) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset, whence)
}

// Readdir implements the Readdir method of the http.File interface
func (f *CustomFile) Readdir(count int) ([]fs.FileInfo, error) {
	if count == 0 {
		return nil, nil
	}
	return nil, fmt.Errorf("CustomFile does not support Readdir")
}

// Stat implements the Stat method of the http.File interface
func (f *CustomFile) Stat() (os.FileInfo, error) {
	return f.FileInfo, nil
}

func (f *CustomFile) Close() error {
	// bytes.Reader doesn't require explicit closing
	return nil
}

type mockFileInfo struct {
	size    int64
	modTime time.Time
	isDir   bool
	name    string
	mode    os.FileMode
}

func (f *mockFileInfo) Size() int64 {
	return f.size
}

func (f *mockFileInfo) ModTime() time.Time {
	return f.modTime
}

func (f *mockFileInfo) IsDir() bool {
	return f.isDir
}

func (f *mockFileInfo) Name() string {
	return f.name
}

func (f *mockFileInfo) Sys() interface{} {
	return nil
}

func (f *mockFileInfo) Mode() os.FileMode {
	// Set a default mode (e.g., regular file)
	return 0644 // Adjust as needed based on your files
}
