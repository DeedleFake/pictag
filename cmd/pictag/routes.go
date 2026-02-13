package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"deedles.dev/pictag/internal/sqlc"
	"deedles.dev/pictag/internal/ui"
)

func initRoutes(h *handler) {
	http.Handle("GET /img/", http.StripPrefix("/img/", http.FileServerFS(h.store.FS())))

	handleAPI("GET /api/list_tags", h.listTags)
	handleAPI("GET /api/list_images", h.listImages)

	http.Handle("GET /", ui.Handler())
}

func handleAPI(route string, h http.HandlerFunc) {
	http.HandleFunc(route, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		h(rw, req)
	}))
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

func (h *handler) listImages(rw http.ResponseWriter, req *http.Request) {
	slog := withRequest(slog.Default(), req)

	err := req.ParseForm()
	if err != nil {
		slog.Error("parse form", "err", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	q := req.Form["q"]
	images, err := sqlc.New(h.db).ImagesByTags(req.Context(), sqlc.ImagesByTagsParams{
		Tags:   q,
		Length: 1,
		Limit:  10,
	})
	if err != nil {
		slog.Error("list images", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(rw).Encode(images)
	if err != nil {
		slog.Error("encode response", "err", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
