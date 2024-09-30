package main

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

func main() {
	var loggingLevel = new(slog.LevelVar)

	// we only want to calculate the working directory once
	// not every time we log
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to determine working directory")
	}

	replacer := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			// remove current working directory and only leave the relative path to the program
			if file, ok := strings.CutPrefix(source.File, wd); ok {
				source.File = file
			}
		}
		return a
	}

	options := &slog.HandlerOptions{
		Level:       loggingLevel,
		ReplaceAttr: replacer,
	}
	// if debug is present, turn on logging level and source for logs
	if os.Getenv("loglevel") == "debug" {
		loggingLevel.Set(slog.LevelDebug)
		options.AddSource = true
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, options))

	// Set the logger for the application
	slog.SetDefault(logger)
	slog.Info("hello gophers")
	slog.Warn("be warned!")
	slog.Error("this is broken")
	slog.Debug("this is a debug")

}
