package ui

import (
	"fmt"
	"strings"

	"container-tui/src/models"
)

// RenderResult formats a result for display.
func RenderResult(result models.Result) string {
	status := strings.ToUpper(string(result.Status))
	statusLine := "Status: " + status
	switch result.Status {
	case models.ResultSuccess:
		statusLine = RenderSuccess(statusLine)
	case models.ResultError:
		statusLine = RenderError(statusLine)
	default:
		statusLine = RenderMuted(statusLine)
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", statusLine))
	if result.Stdout != "" {
		builder.WriteString("\nStdout:\n")
		builder.WriteString(result.Stdout)
	}
	if result.Stderr != "" {
		builder.WriteString("\nStderr:\n")
		builder.WriteString(result.Stderr)
	}
	return builder.String()
}
