package ui

//go:generate bun run --bun build

import (
	"embed"
	"io/fs"
	"net/http"
)

var (
	//go:embed assets
	assets embed.FS
	sub    = mustSub(assets, "assets")
)

func mustSub(root fs.FS, dir string) fs.FS {
	r, err := fs.Sub(root, dir)
	if err != nil {
		panic(err)
	}
	return r
}

func FS() fs.FS {
	return sub
}

func handleIndex(rw http.ResponseWriter, req *http.Request) {
	http.ServeFileFS(rw, req, assets, "assets/index.html")
}

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.FileServerFS(assets))
	mux.HandleFunc("GET /", handleIndex)
	return mux
}
