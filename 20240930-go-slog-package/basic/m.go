package main

import "log/slog"

func main() {
	slog.Info("hello gophers")
	slog.Warn("be warned!")
	slog.Error("this is broken")
	slog.Debug("show some debugging output")
}
