package launchy

import (
	"fmt"
	"io"
	"os"
)

// Logger holds information about how Launchy should handle logging information
type Logger struct {
	Stdout  io.Writer
	Stderr  io.Writer
	Verbose bool
}

// NewLogger creates a new Launchy logger.
func NewLogger(outputWriter, errorWriter io.Writer, verbose bool) *Logger {
	if outputWriter == nil {
		outputWriter = os.Stdout
	}

	if errorWriter == nil {
		errorWriter = os.Stderr
	}

	return &Logger{
		Stdout:  outputWriter,
		Stderr:  errorWriter,
		Verbose: verbose,
	}
}

// VerbosePrintf prints a formatted message to the provided output
// io.Writer when the Verbose flag is set. By default, if this is not provided
// when creating a logger instance, os.Stdout will be used.
func (l *Logger) VerbosePrintf(format string, v ...interface{}) {
	if l.Verbose {
		fmt.Fprintf(l.Stdout, "[launchy] "+format+"\n", v)
	}
}

// VerboseErrorPrintf prints a formatted message to the provided output
// io.Writer when the Verbose flag is set. By default, if this is not provided
// when creating a logger instance, os.Stderr will be used.
func (l *Logger) VerboseErrorPrintf(format string, v ...interface{}) {
	if l.Verbose {
		fmt.Fprintf(l.Stderr, "[launchy] "+format+"\n", v)
	}
}

// Printf prints a formatted message to the provided output
// io.Writer. By default, if this is not provided when creating a logger
// instance, os.Stderr will be used.
func (l *Logger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(l.Stdout, format+"\n", v)
}

// ErrorPrintf prints a formatted message to the provided output
// io.Writer. By default, if this is not provided when creating a logger
// instance, os.Stderr will be used.
func (l *Logger) ErrorPrintf(format string, v ...interface{}) {
	fmt.Fprintf(l.Stderr, format+"\n", v)
}
