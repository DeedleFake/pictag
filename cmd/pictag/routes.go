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

	http.HandleFunc("GET /api/list_tags", h.listTags)

	http.HandleFunc("GET /", h.index)
}

func (h *handler) listTags(rw http.ResponseWriter, req *http.Request) {
	slog := withRequest(slog.Default(), req)
	rw.Header().Set("Content-Type", "application/json")

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

	images, err := sqlc.New(h.db).ImagesByTags(req.Context(), sqlc.ImagesByTagsParams{
		Tags:   []string{"test"},
		Length: 1,
		Limit:  10,
	})
	if err != nil {
		slog.Error("list images", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	names := func(yield func(string) bool) {
		for _, img := range images {
			if !yield(img.ID) {
				return
			}
		}
	}

	err = ui.Index(names).Render(req.Context(), rw)
	if err != nil {
		slog.Error("render component", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
