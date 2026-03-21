package logger

import (
	"fmt"
	"os"
	"strings"
)

// -----------------------------------------------------
//
//   Logger is responsible for logging messages to the terminal.
//   It writes exclusively to stderr — stdout is reserved for the cd command
//   emitted to the shell for path navigation.
//
// -----------------------------------------------------

type Logger struct{}

var LOGGER *Logger

// ----------------------------------
//
//	Initializes the global LOGGER instance for use across the application.
//
// ----------------------------------
func InitLogger() {
	LOGGER = NewLogger()
}

// ----------------------------------
//
//	Creates and returns a new Logger instance.
//
// ----------------------------------
func NewLogger() *Logger {
	return &Logger{}
}

// ----------------------------------
//
//	Writes all messages to stderr, joined by newlines.
//
// ----------------------------------
func (l *Logger) LogToTerminal(message []string) {
	fmt.Fprintln(os.Stderr, strings.Join(message, "\n"))
}
