package logger

import (
	"fmt"
	"os"
	"strings"
)

// -----------------------------------------------------
//
//   Logger was resposible for logging messages to the console and files.
//   It will output using stderr as stdout was solely for the the delegate of command to the shell for path navigation and execution.
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
