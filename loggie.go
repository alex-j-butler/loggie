package loggie

import (
	"fmt"
	"os"
	"time"
)

type LogLevel uint16
type logOutput struct {
	name   string
	output func(level LogLevel, logStr string) error
}
type Logger struct {
	Name      string
	ErrLogger ErrorLogger
	outputs   []*logOutput
}
type ErrorLogger func(logStr string)

const (
	Fatal LogLevel = iota
	Error
	Warn
	Info
	Debug
)

var (
	LevelToString = map[LogLevel]string{
		Fatal: "FATAL",
		Error: "ERROR",
		Warn:  "WARN",
		Info:  "INFO",
		Debug: "DEBUG",
	}
)

func DefaultErrorLogger() ErrorLogger {
	return func(logStr string) {
		fmt.Fprintf(os.Stderr, "%s\n", logStr)
	}
}

func StdLogger() *logOutput {
	return &logOutput{
		name: "StdLogger",
		output: func(level LogLevel, logStr string) error {
			logLine := fmt.Sprintf("[%s %s] %s\n", time.Now().Format("02/Jan/2006:15:04:05 -0700"), LevelToString[level], logStr)

			if level <= Warn {
				fmt.Fprintf(os.Stderr, logLine)
			} else {
				fmt.Fprintf(os.Stdout, logLine)
			}
			return nil
		},
	}
}

func NewNamedLogger(name string, outputs ...*logOutput) *Logger {
	return &Logger{
		Name:      name,
		ErrLogger: DefaultErrorLogger(),
		outputs:   outputs,
	}
}

func (logger *Logger) rawLog(level LogLevel, logStr string) (map[string]error, bool) {
	hasError := false // Whether any of the logger outputs have failed.
	errorMap := make(map[string]error)
	for _, output := range logger.outputs {
		if err := output.output(level, logStr); err != nil {
			hasError = true
			errorMap[output.name] = err
		}
	}

	return errorMap, hasError
}

func (logger *Logger) logErrors(errMap map[string]error) {
	for name, err := range errMap {
		logger.ErrLogger(fmt.Sprintf("Error logging to %s: %v", name, err))
	}
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	errMap, fail := logger.rawLog(Fatal, fmt.Sprintf(format, args...))
	if fail {
		logger.logErrors(errMap)
	}
	os.Exit(1) // Exit the application, fatal error.
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	errMap, fail := logger.rawLog(Error, fmt.Sprintf(format, args...))
	if fail {
		logger.logErrors(errMap)
	}
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	errMap, fail := logger.rawLog(Warn, fmt.Sprintf(format, args...))
	if fail {
		logger.logErrors(errMap)
	}
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	errMap, fail := logger.rawLog(Info, fmt.Sprintf(format, args...))
	if fail {
		logger.logErrors(errMap)
	}
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	errMap, fail := logger.rawLog(Debug, fmt.Sprintf(format, args...))
	if fail {
		logger.logErrors(errMap)
	}
}
