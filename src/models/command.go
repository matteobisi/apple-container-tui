package models

import (
	"strconv"
	"strings"
)

// Command represents an executable and its arguments.
type Command struct {
	Executable string
	Args       []string
}

// String formats the command for previews.
func (c Command) String() string {
	parts := make([]string, 0, len(c.Args)+1)
	if c.Executable != "" {
		parts = append(parts, c.Executable)
	}
	for _, arg := range c.Args {
		parts = append(parts, quoteIfNeeded(arg))
	}
	return strings.Join(parts, " ")
}

func quoteIfNeeded(value string) string {
	if value == "" || strings.ContainsAny(value, " \t\n\"'") {
		return strconv.Quote(value)
	}
	return value
}
