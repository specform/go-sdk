package main

import (
	"log/slog"
	"os"
)

func NewLogger(verbose bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if verbose {
		opts.Level = slog.LevelDebug
	}
	return slog.New(slog.NewJSONHandler(os.Stderr, opts))
}
