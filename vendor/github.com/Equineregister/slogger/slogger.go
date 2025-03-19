package slogger

import (
	"fmt"
	"log/slog"
	"os"
)

// OptFunc is a type function definition to allow optional field updates
// on the local config.
type OptFunc func(c *config)

// WithLogLevel exposes the option to overwrite the default log Level
func WithLogLevel(level slog.Level) OptFunc {
	return func(c *config) {
		c.Loglevel = &level
	}
}

// WithSource exposes the option to overwrite the default value (true)
// of `IncludeSource`.
func WithSource(source bool) OptFunc {
	return func(c *config) {
		c.IncludeSource = source
	}
}

// WithStyle exposes the option to overwrite the default logger style
// that is associated with the `Env`.
func WithStyle(style Style) OptFunc {
	return func(c *config) {
		c.Style = style
	}
}

type config struct {
	IncludeSource bool
	Style         Style
	Loglevel      *slog.Level
}

// defaultConfig returns a `config` populated with the default options.
func defaultConfig() *config {
	return &config{
		IncludeSource: true,
		Style:         Style(""),
		Loglevel:      nil,
	}
}

var envOpts = map[Env]struct {
	Style    Style
	Loglevel slog.Level
}{
	EnvDev: {
		Style:    StyleJSON,
		Loglevel: slog.LevelDebug,
	},
	EnvUat: {
		Style:    StyleJSON,
		Loglevel: slog.LevelInfo,
	},
	EnvProd: {
		Style:    StyleJSON,
		Loglevel: slog.LevelWarn,
	},
	EnvLocal: {
		Style:    StyleLocal,
		Loglevel: slog.LevelDebug,
	},
}

// NewHandler returns a new `slog.Handler` based on the provided `Env`.
// For example, if the `Env` is `EnvDev` NewHandler will set the log level
// to `slog.LevelInfo` and the style to JSON.
// These values can be overwritten by passing an optional `OptFunc`.
//
// Basic usage (Returns a new `slog.Logger` ready for DEV env):
//
// ```go
//
//	slogHandler := slogger.NewHandler(slogger.EnvDev)
//	logger := slog.New(slog.Handler)
//
// ```
//
// Usage with config (Overriding the style):
// ```go
//
//	dev := slogger.NewHandler(slogger.EnvDev, slogger.WithStyle(slogger.StyleJSON))
//	logger := slog.New(slog.Handler)
//
// ```
func NewHandler(env Env, opts ...OptFunc) (slog.Handler, error) {
	if !env.isValid() {
		return nil, fmt.Errorf("invalid env: %s", env.String())
	}

	config := defaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	handlerOpts := &slog.HandlerOptions{
		AddSource: config.IncludeSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove the time from the log output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}

	// Set the log level if it is provided.
	shouldOverwriteLogLovel := config.Loglevel != nil
	if shouldOverwriteLogLovel {
		handlerOpts.Level = config.Loglevel.Level()
	}

	// Get the style and log level for the provided env.
	envOpts, ok := envOpts[env]
	if !ok {
		return nil, fmt.Errorf("unknown env provided %q", env)
	}

	// Only use the env log level if the config log level is not provided.
	if !shouldOverwriteLogLovel {
		handlerOpts.Level = envOpts.Loglevel.Level()
	}

	// If a `config.Style` is populated, override the logger style.
	if config.Style.isValid() {
		return handlerFromStyle(handlerOpts, config.Style)
	}

	return handlerFromStyle(handlerOpts, envOpts.Style)
}

func handlerFromStyle(handlerOpts *slog.HandlerOptions, style Style) (slog.Handler, error) {
	switch style {
	case StyleJSON:
		return slog.NewJSONHandler(os.Stdout, handlerOpts), nil

	case StyleText:
		return slog.NewTextHandler(os.Stdout, handlerOpts), nil

	case StyleLocal:
		return NewLocalHandler(os.Stdout, handlerOpts.Level), nil

	default:
		return nil, fmt.Errorf("unknown log style provided %q", style)
	}
}
