package main

import "alex-j-butler.com/loggie"

func main() {
	logger := loggie.NewNamedLogger("example", loggie.StdLogger())
	logger.Debugf("example debug")
	logger.Infof("example info")
	logger.Warnf("example warn")
	logger.Errorf("example error")
	logger.Fatalf("example fatal")
}
