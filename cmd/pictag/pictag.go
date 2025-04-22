package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"deedles.dev/pictag/store"
	"github.com/adrg/xdg"
	_ "modernc.org/sqlite"
)

func initRoutes(data string) {
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
	defer store.Close()
	slog.Info("store opened", "path", data)
}

func main() {
	data := flag.String("data", filepath.Join(xdg.StateHome, "pictag"), "directory to store data in")
	addr := flag.String("addr", "localhost:5050", "address to listen on for HTTP")
	flag.Parse()

	initRoutes(*data)

	slog.Info("listening for HTTP", "address", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		slog.Error("listen and serve HTTP", "err", err)
		os.Exit(1)
	}
}
