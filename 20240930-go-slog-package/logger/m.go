package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Info("hello gophers")
	logger.Warn("be warned!")
	logger.Error("this is broken")
}
