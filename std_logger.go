package loggie

import (
	"fmt"
	"os"
	"time"
)

// StdLogger is a logging output that prints log lines to stdout and stderr.
type StdLogger struct {
	// Whether the output should combine stderr and stdout and directly output to stdout.
	// True = combine stderr and stdout
	// DEPRECATED: Doesn't make sense for a stdout/stderr logger.
	Combined bool
}

func NewStdLogger() *StdLogger {
	return &StdLogger{
		Combined: false,
	}
}

func (l StdLogger) GetName() string {
	return "StdLogger"
}

func (l StdLogger) Output(logLevel LogLevel, logString string) error {
	logLine := fmt.Sprintf("[%s %s] %s\n", time.Now().Format("02/Jan/2006:15:04:05 -0700"), LevelToString[logLevel], logString)

	if logLevel <= Warn {
		fmt.Fprintf(os.Stderr, logLine)
	} else {
		fmt.Fprintf(os.Stdout, logLine)
	}
	return nil
}
