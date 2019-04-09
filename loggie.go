package loggie

import (
	"fmt"
	"os"
	"time"
)

type LoggingOutput interface {
	GetName() string
	Output(level LogLevel, logString string) error
}

type LogFormatter interface {
	GetName() string
	Format(t time.Time, level LogLevel, logString string) string
}

// StdLogger is a logging output that prints log lines to stdout and stderr.
type StdLogger struct {
	// Whether the output should combine stderr and stdout and directly output to stdout.
	// True = combine stderr and stdout
	// DEPRECATED: Doesn't make sense for a stdout/stderr logger.
	Combined bool
}

type FileLogger struct {
	Stdout *os.File
	Stderr *os.File
}

type LogLevel uint16
type LogOutput struct {
	Name   string
	Output func(level LogLevel, logStr string) error
}
type Logger struct {
	Name      string
	Level     LogLevel
	ErrLogger ErrorLogger
	outputs   []*LogOutput
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

func NewStdLogger() *StdLogger {
	return &StdLogger{
		Combined: false,
	}
}

func CombinedFileLogger(file *os.File) *LogOutput {
	return &LogOutput{
		Name: "CombinedFileLogger",
		Output: func(level LogLevel, logStr string) error {
			logLine := fmt.Sprintf("[%s %s] %s\n", time.Now().Format("02/Jan/2006:15:04:05 -0700"), LevelToString[level], logStr)
			fmt.Fprintf(file, logLine)

			return nil
		},
	}
}

func NewNamedLogger(name string, logLevel LogLevel, outputs ...*LogOutput) *Logger {
	return &Logger{
		Name:      name,
		Level:     logLevel,
		ErrLogger: DefaultErrorLogger(),
		outputs:   outputs,
	}
}

func NewLogger(outputs ...*LogOutput) *Logger {
	return &Logger{
		ErrLogger: DefaultErrorLogger(),
		outputs:   outputs,
	}
}

func (logger *Logger) rawLog(level LogLevel, logStr string) (map[string]error, bool) {
	if level > logger.Level {
		return nil, false
	}

	hasError := false // Whether any of the logger outputs have failed.
	errorMap := make(map[string]error)
	for _, output := range logger.outputs {
		if err := output.Output(level, logStr); err != nil {
			hasError = true
			errorMap[output.Name] = err
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
