package main

import (
	"errors"
	"fmt"
	"strings"
)

// Color is a simple enum for various colors.
type Color int8

// All available colors
const (
	_ Color = iota
	Red
	Blue
	Green
)

// ParseColor takes a color name and returns enum value.
func ParseColor(colorName string) (Color, error) {
	switch strings.ToLower(colorName) {
	case "red":
		return Red, nil
	case "blue":
		return Blue, nil
	case "green":
		return Green, nil
	default:
		return 0, fmt.Errorf("unrecogized color name %q", colorName)
	}
}

// String returns the color name with the first letter capitalized.
func (c Color) String() string {
	switch c {
	case Red:
		return "Red"
	case Blue:
		return "Blue"
	case Green:
		return "Green"
	default:
		return ""
	}
}

// MarshalText implements the TextMarshaler interface.
func (c Color) MarshalText() ([]byte, error) {
	s := strings.ToLower(c.String())
	if s == "" {
		return nil, errors.New("invalid color")
	}
	return []byte(s), nil
}

// UnmarshalText implements the TextUnmarshaler interface.
func (c *Color) UnmarshalText(buf []byte) error {
	parsed, err := ParseColor(string(buf))
	if err != nil {
		return err
	}
	*c = parsed
	return nil
}
