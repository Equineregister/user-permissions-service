package slogger

import (
	"errors"
)

const (
	StyleText  Style = "text"
	StyleJSON  Style = "json"
	StyleLocal Style = "local"
)

type Style string

func (s *Style) isValid() bool {
	switch *s {
	case StyleText, StyleJSON, StyleLocal:
		return true
	}

	return false
}

func (s *Style) Set(value string) error {
	style := Style(value)
	if !style.isValid() {
		return errors.New(`style must be one of "text", "json", or "local"`)
	}

	*s = style

	return nil
}

func (s Style) String() string {
	return string(s)
}
