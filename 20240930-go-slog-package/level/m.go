package main

import (
	"log/slog"
	"os"
)

func main() {
	// Pass the loggingLevel to the new logger being created so we can change it later at any time
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger.Info("hello gophers")
	logger.Warn("be warned!")
	logger.Error("this is broken")
	logger.Debug("be warned!")
}
