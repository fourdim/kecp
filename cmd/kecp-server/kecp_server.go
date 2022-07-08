package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/fourdim/kecp/router"
	"github.com/pelletier/go-toml/v2"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var App struct {
	Server struct {
		Debug          bool
		TLS            bool
		Host           string
		AllowedOrigins []string `toml:"allowed_origins"`
	}
}

func main() {
	b, err := os.ReadFile("config.toml")
	if err != nil {
		log.Panicln(err)
	}

	toml.Unmarshal(b, &App)

	kecpApiServerRouter := chi.NewRouter()

	kecpApiServerRouter.Route("/api", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: App.Server.AllowedOrigins,
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Upgrade", "Connection", "Sec-WebSocket-Key", "Sec-WebSocket-Protocol", "Sec-WebSocket-Version", "Sec-WebSocket-Extensions"},
			ExposedHeaders:   []string{"Sec-WebSocket-Accept", "Sec-WebSocket-Protocol", "Sec-WebSocket-Extensions"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Mount("/kecp", router.SetupKecpChiRouter())
	})

	kecpApiServerRouter.NotFound(ServeRoot("/", "./app/dist"))

	if App.Server.Debug || !App.Server.TLS {
		http.ListenAndServe(":8090", kecpApiServerRouter)
	} else {
		certmagic.HTTPS([]string{App.Server.Host}, kecpApiServerRouter)
	}

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
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

func ServeRoot(urlPrefix, root string) http.HandlerFunc {
	return Serve(urlPrefix, root, LocalFile(root, false))
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, root string, fs ServeFileSystem) http.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if fs.Exists(urlPrefix, r.URL.Path) {
			fileserver.ServeHTTP(w, r)
		} else {
			f, err := os.Open(path.Join(root, INDEX))
			if err != nil {
				return
			}
			http.ServeContent(w, r, INDEX, time.Now(), f)
		}
	}
}
