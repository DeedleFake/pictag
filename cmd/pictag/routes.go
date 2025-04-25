package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"deedles.dev/pictag/internal/assets"
	"deedles.dev/pictag/internal/sqlc"
	"deedles.dev/pictag/internal/ui"
)

func initRoutes(h *handler) {
	http.Handle("GET /assets/", http.StripPrefix("/assets/", assets.Handler()))
	http.Handle("GET /img/", http.StripPrefix("/img/", http.FileServerFS(h.store.FS())))

	http.HandleFunc("GET /ui/list_tags", h.listTags)

	http.HandleFunc("GET /", h.index)
}

func (h *handler) listTags(rw http.ResponseWriter, req *http.Request) {
	slog := withRequest(slog.Default(), req)

	q := req.FormValue("q")
	tags, err := sqlc.New(h.db).SearchTags(req.Context(), sqlc.SearchTagsParams{
		Name:  q + "%",
		Limit: 10,
	})
	if err != nil {
		slog.Error("search tags", "query", q, "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(tags)
	if err != nil {
		slog.Error("encode response", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) index(rw http.ResponseWriter, req *http.Request) {
	slog := withRequest(slog.Default(), req)

	err := ui.Index().Render(req.Context(), rw)
	if err != nil {
		slog.Error("render component", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
