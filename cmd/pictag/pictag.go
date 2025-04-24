package main

import (
	"context"
	"database/sql"
	"flag"
	"image"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"deedles.dev/pictag/internal/assets"
	"deedles.dev/pictag/internal/sqlc"
	"deedles.dev/pictag/store"
	"github.com/adrg/xdg"
	_ "modernc.org/sqlite"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "deedles.dev/ximage/xcursor"
	_ "github.com/HugoSmits86/nativewebp"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

type handler struct {
	store *store.Store
	db    *sql.DB
}

func initHandler(data string) *handler {
	err := os.MkdirAll(data, 0755)
	if err != nil {
		slog.Error("create data directory", "path", data, "err", err)
		os.Exit(1)
	}

	store, err := store.Open(data)
	if err != nil {
		slog.Error("open store", "err", err)
		os.Exit(1)
	}
	slog.Info("store opened", "path", data)

	dbpath := filepath.Join(data, "data.db")
	db, err := sql.Open("sqlite", dbpath)
	if err != nil {
		slog.Error("open database", "path", dbpath, "err", err)
		os.Exit(1)
	}
	slog.Info("database opened", "path", dbpath)

	err = sqlc.Migrate(context.TODO(), db)
	if err != nil {
		slog.Error("migrate database", "err", err)
		os.Exit(1)
	}

	return &handler{
		store: store,
		db:    db,
	}
}

func initRoutes(h *handler) {
	http.Handle("GET /assets/", http.StripPrefix("/assets/", assets.Handler()))
	http.Handle("GET /img/", http.StripPrefix("/img/", http.FileServerFS(h.store.FS())))

	http.HandleFunc("POST /test", func(rw http.ResponseWriter, req *http.Request) {
		file, _, err := req.FormFile("image")
		if err != nil {
			slog.Error("parse image form data", "err", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			slog.Error("decode image", "err", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		name, err := h.store.Store(img)
		if err != nil {
			slog.Error("store image", "err", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.WriteString(rw, name)
	})
}

func main() {
	data := flag.String("data", filepath.Join(xdg.StateHome, "pictag"), "directory to store data in")
	addr := flag.String("addr", "localhost:5050", "address to listen on for HTTP")
	flag.Parse()

	h := initHandler(*data)
	initRoutes(h)

	slog.Info("listening for HTTP", "address", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		slog.Error("listen and serve HTTP", "err", err)
		os.Exit(1)
	}
}
