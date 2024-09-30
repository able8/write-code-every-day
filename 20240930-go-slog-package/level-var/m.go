package main

import (
	"log/slog"
	"os"
)

func main() {
	// create a logging level variable
	// the level is Info by default
	var loggingLevel = new(slog.LevelVar)

	// Pass the loggingLevel to the new logger being created so we can change it later at any time
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: loggingLevel}))

	logger.Info("hello gophers")
	logger.Warn("be warned!")
	logger.Error("this is broken")
	logger.Debug("be warned!")

	logger.Info("---------")

	loggingLevel.Set(slog.LevelDebug)

	logger.Info("hello gophers")
	logger.Warn("be warned!")
	logger.Error("this is broken")
	logger.Debug("be warned!")

	logger.Info("---------")

	loggingLevel.Set(slog.LevelError)

	logger.Info("hello gophers")
	logger.Warn("be warned!")
	logger.Error("this is broken")
	logger.Debug("be warned!")
}
