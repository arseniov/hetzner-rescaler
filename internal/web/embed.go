// Package web embeds the SvelteKit static SPA build output and serves it
// as a fallback handler. Any request that does not match a real file in
// the embed (and that isn't an /api/* call — handled upstream) is served
// index.html so the SPA's client-side router can take over.
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:build
var buildFS embed.FS

// indexHTML is read once at startup. It is served for any request that
// does not match a real file in the embed.
var indexHTML []byte

func init() {
	b, err := buildFS.ReadFile("build/index.html")
	if err != nil {
		// index.html may not exist yet during early development; fall
		// back to a placeholder so the binary still starts. Production
		// builds always include index.html (SvelteKit adapter-static
		// emits it).
		indexHTML = []byte("<!doctype html><html><body>rescaler web UI: " +
			"run `bun run build` inside web/</body></html>")
		return
	}
	indexHTML = b
}

// SPAHandler returns an http.Handler that serves the embedded build.
// Files that exist in the embed are served with their content type
// inferred from the extension. Anything else falls back to index.html.
func SPAHandler() http.Handler {
	sub, err := fs.Sub(buildFS, "build")
	if err != nil {
		// fs.Sub only fails if the root is invalid; recover with a
		// handler that always serves the placeholder.
		return fallbackOnly()
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			writeIndex(w)
			return
		}
		// If the file exists in the embed, serve it.
		if _, err := fs.Stat(sub, path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}
		// Otherwise, fall back to index.html.
		writeIndex(w)
	})
}

func writeIndex(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(indexHTML)
}

func fallbackOnly() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeIndex(w)
	})
}
