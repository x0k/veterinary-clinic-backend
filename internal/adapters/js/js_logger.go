//go:build js && wasm

package js_adapters

import (
	"context"
	"log/slog"
)

const Disabled slog.Level = -8

type ConsoleLoggerHandler struct {
	level slog.Level
	attrs []slog.Attr
}

func NewConsoleLoggerHandler(level slog.Level) *ConsoleLoggerHandler {
	return &ConsoleLoggerHandler{
		level: level,
	}
}

func (h *ConsoleLoggerHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *ConsoleLoggerHandler) Handle(_ context.Context, r slog.Record) error {
	logMethod := "debug"
	switch r.Level {
	case slog.LevelDebug:
		logMethod = "debug"
	case slog.LevelInfo:
		logMethod = "info"
	case slog.LevelWarn:
		logMethod = "warn"
	case slog.LevelError:
		logMethod = "error"
	}

	fields := make(map[string]any, r.NumAttrs()+len(h.attrs))

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	timeLabel := r.Time.Format("[15:04:05.000]")

	Console.Get(logMethod).Invoke(timeLabel, r.Message, fields)
	return nil
}

func (h *ConsoleLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ConsoleLoggerHandler{
		level: h.level,
		attrs: append(h.attrs, attrs...),
	}
}

// TODO: add WithGroup
func (h *ConsoleLoggerHandler) WithGroup(name string) slog.Handler {
	return h
}
