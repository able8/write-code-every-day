package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{AddSource: true}))
	// Set the logger for the application
	slog.SetDefault(logger)

	slog.Info("hello gophers")
	slog.Warn("be warned!")
	slog.Error("this is broken")

}
