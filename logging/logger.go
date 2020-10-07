package logging

import (
	"fmt"
	"time"
)

const (
	logLevelQuiet = iota
	logLevelError
	logLevelInfo
	logLevelDebug
)

var logLevel = logLevelInfo

// Initialize initializes logging.
func Initialize(level int) {
	if level >= logLevelQuiet && level <= logLevelDebug {
		logLevel = level
	}
}

// Log set log.
func Log(format string, attr ...interface{}) {
	if logLevel >= logLevelInfo {
		fmt.Printf(prependTime("INFO   "+format+"\n", attr...))
	}
	return
}

// Debug set log.
func Debug(format string, attr ...interface{}) {
	if logLevel >= logLevelDebug {
		fmt.Printf(prependTime("DEBUG  "+format+"\n", attr...))
	}
	return
}

// Error set log.
func Error(format string, attr ...interface{}) {
	if logLevel >= logLevelError {
		fmt.Printf(prependTime("ERROR  "+format+"\n", attr...))
	}
	return
}

// Fatal set log.
func Fatal(format string, attr ...interface{}) {
	panic(prependTime(format+"\n", attr...))
}

func prependTime(format string, attr ...interface{}) string {
	t := fmt.Sprintf("%v", time.Now().Format(time.RFC3339))
	return t + "  " + fmt.Sprintf(format, attr...)
}
