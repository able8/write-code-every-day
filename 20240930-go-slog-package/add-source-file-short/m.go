package main

import (
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	replacer := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	logger := slog.New(slog.NewJSONHandler(
		os.Stderr, &slog.HandlerOptions{
			AddSource:   true,
			ReplaceAttr: replacer,
		}))
	// Set the logger for the application
	slog.SetDefault(logger)

	slog.Info("hello gophers")
	slog.Warn("be warned!")
	slog.Error("this is broken")

}
