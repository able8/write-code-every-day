package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	log.Println("program started")
	r, err := http.Get("https://www.gopherguides.com")
	if err != nil {
		slog.Error("error retrieving site", slog.String("err", err.Error()))
	}
	slog.Info("success",
		slog.Group(
			"request",
			slog.String("method", r.Request.Method),
			slog.String("url", r.Request.URL.String()),
		))
}
