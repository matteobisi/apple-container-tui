package services

import "container-tui/src/models"

// LoggingExecutor wraps an executor with log writing.
type LoggingExecutor struct {
	delegate CommandExecutor
	writer   *LogWriter
	dryRun   bool
}

// NewLoggingExecutor builds a logging executor.
func NewLoggingExecutor(delegate CommandExecutor, writer *LogWriter, dryRun bool) *LoggingExecutor {
	return &LoggingExecutor{delegate: delegate, writer: writer, dryRun: dryRun}
}

// Execute runs the command and writes a log entry.
func (l *LoggingExecutor) Execute(cmd models.Command) (models.Result, error) {
	result, err := l.delegate.Execute(cmd)
	if l.writer != nil {
		_ = l.writer.Write(BuildLogEntry(cmd, result, l.dryRun))
	}
	return result, err
}
