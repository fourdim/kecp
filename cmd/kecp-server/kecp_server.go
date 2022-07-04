package main

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fourdim/kecp/router"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	kecpApiServerRouter := chi.NewRouter()

	kecpApiServerRouter.Route("/api", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Upgrade", "Connection", "Sec-WebSocket-Key", "Sec-WebSocket-Protocol", "Sec-WebSocket-Version", "Sec-WebSocket-Extensions"},
			ExposedHeaders:   []string{"Sec-WebSocket-Accept", "Sec-WebSocket-Protocol", "Sec-WebSocket-Extensions"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Mount("/kecp", router.SetupKecpChiRouter())
	})

	kecpApiServerRouter.NotFound(ServeRoot("/", "./tmp"))

	http.ListenAndServe(":8090", kecpApiServerRouter)
}

const INDEX = "index.html"

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: http.Dir(root),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, INDEX)
				stats, err := os.Stat(index)
				if err != nil && !stats.IsDir() {
					return false
				}
			}
		}
		return true
	}
	return false
}

func ServeRoot(urlPrefix, root string) http.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) http.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if fs.Exists(urlPrefix, r.URL.Path) {
			fileserver.ServeHTTP(w, r)
		}
	}
}
