package slogger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"sync"
)

const LocalTimeFormat = "15:04:05 MST"

// LocalHandler implements as slog.Handler that format logs in a digestible
// format for local development. It is not recommended to use this handler within
// any environment other than local development.
type LocalHandler struct {
	level slog.Leveler
	group string
	attrs []slog.Attr
	mu    *sync.Mutex
	w     io.Writer
}

// NewLocalHandler returns a pointer to a new LocalHandler. The handler will
// write to the provided io.Writer and will only log records that are greater
// than or equal to the provided level.
func NewLocalHandler(w io.Writer, level slog.Leveler) *LocalHandler {
	if level == nil || reflect.TypeOf(level).Kind() == reflect.Ptr && reflect.ValueOf(level).IsNil() {
		level = &slog.LevelVar{}
	}

	return &LocalHandler{
		level: level,
		mu:    new(sync.Mutex),
		w:     w,
	}
}

func (h *LocalHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *LocalHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var attrs string
	for _, a := range h.attrs {
		if !a.Equal(slog.Attr{}) {
			attrs += "  "

			if h.group != "" {
				attrs += h.group + "."
			}

			attrs += a.Key + ": " + a.Value.String() + "\n"
		}
	}

	r.Attrs(func(a slog.Attr) bool {
		if !a.Equal(slog.Attr{}) {
			attrs += "  "

			if h.group != "" {
				attrs += h.group + "."
			}

			attrs += a.Key + ": " + a.Value.String() + "\n"
		}

		return true
	})

	attrs = strings.TrimRight(attrs, "\n")

	var newlines string
	if attrs != "" {
		newlines = "\n\n"
	}

	fmt.Fprintf(h.w, "[%v] %v\n%v%v", r.Time.Format(LocalTimeFormat), r.Message, attrs, newlines)

	return nil
}

func (h *LocalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LocalHandler{
		level: h.level,
		group: h.group,
		attrs: append(h.attrs, attrs...),
		mu:    h.mu,
		w:     h.w,
	}
}

func (h *LocalHandler) WithGroup(name string) slog.Handler {
	return &LocalHandler{
		level: h.level,
		group: strings.TrimSuffix(name+"."+h.group, "."),
		attrs: h.attrs,
		mu:    h.mu,
		w:     h.w,
	}
}
