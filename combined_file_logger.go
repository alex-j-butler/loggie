package loggie

import (
	"fmt"
	"os"
	"time"
)

// CombinedFileLogger is a logging output that prints log lines to a single file.
type CombinedFileLogger struct {
	File *os.File
}

func (l CombinedFileLogger) GetName() string {
	return "CombinedFileLogger"
}

func (l CombinedFileLogger) Output(logLevel LogLevel, logString string) error {
	logLine := fmt.Sprintf("[%s %s] %s\n", time.Now().Format("02/Jan/2006:15:04:05 -0700"), LevelToString[logLevel], logString)
	fmt.Fprintf(l.File, logLine)

	return nil
}

func NewCombinedFileLogger(file *os.File) *CombinedFileLogger {
	return &CombinedFileLogger{
		File: file,
	}
}
