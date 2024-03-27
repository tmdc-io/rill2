package web

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/NYTimes/gziphandler"
)

//go:embed all:embed
var distFS embed.FS

// Handler serves an web-local UI.
func StaticHandler() http.Handler {
	uiAssetFS := newUIAssetFS()
	return gziphandler.GzipHandler(http.FileServer(uiAssetFS))
}

func DynamicHandler() http.Handler {
	baseDir := os.Getenv("BASE_PATH")
	return gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimPrefix(path, baseDir)
		if path == "" {
			path = "index.html" // Serve index.html if root is requested
		}
		l := "embed/dist/" + path
		//fmt.Printf("3. path=%s; location=%s\n", r.URL.Path, l)
		// Check if the requested path exists in the embedded filesystem
		file, err := distFS.ReadFile(l)
		if err != nil {
			// fallback to index.html
			file, err = distFS.ReadFile("embed/dist/index.html")
			if err != nil {
				fmt.Printf("embed/dist/index.html missing! err=%+v", err)
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
		}

		// Replace "dataos-basepath" with the value of BASE_PATH
		fileContent := strings.ReplaceAll(string(file), "dataos-basepath", baseDir)
		if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, ".html") {
			w.Header().Set("Content-Type", "text/html")
		}
		_, err = w.Write([]byte(fileContent))
		if nil != err {
			fmt.Printf("Failed to write response! err=%+v\n", err)
		}
	}))
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

func (spaFS *SPARoutingFS) Open(name string) (http.File, error) {
	file, err := spaFS.FileSystem.Open(name)
	if err == nil {
		return file, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		file, err = spaFS.FileSystem.Open("index.html")
		return file, err
	}

	return nil, err
}
