package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger creates new logger without work methods.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler is handler in slog.Handler.
type DiscardHandler struct{}

// NewDiscardHandler creates new handler for logger without work methods.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Enabled reports whether the handler handles records at the given level.
func (h *DiscardHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

// Handle handles the Record.
func (h *DiscardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup - A Handler should treat WithGroup as starting a Group of Attrs that ends
// at the end of the log event.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}
