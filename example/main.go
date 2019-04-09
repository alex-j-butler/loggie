package main

import (
	"os"

	"alex-j-butler.com/loggie"
)

func main() {
	combined, err := os.OpenFile("combined.log", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}

	logger := loggie.NewNamedLogger("example", loggie.Debug, loggie.NewStdLogger(), loggie.NewCombinedFileLogger(combined))
	logger.Debugf("example debug")
	logger.Infof("example info")
	logger.Warnf("example warn")
	logger.Errorf("example error")
	logger.Fatalf("example fatal")
}
