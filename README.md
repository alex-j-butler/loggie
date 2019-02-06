# Loggie

A simple logging Golang package.

## Usage

```golang
logger = loggie.NewNamedLogger("example-logger", loggie.StdLogger())
logger.Debugf("debug statement")
logger.Infof("info statement")
logger.Warningf("warning statement")
logger.Errorf("error statement")
logger.Fatalf("fatal statement, program terminated")
```
