package logging

import "fmt"

// Log set log.
func Log(format string, attr ...interface{}) {
	fmt.Printf(format+"\n", attr...)
	return
}

// Debug set log.
func Debug(format string, attr ...interface{}) {
	fmt.Printf(format+"\n", attr...)
	return
}

// Error set log.
func Error(format string, attr ...interface{}) {
	fmt.Printf(format+"\n", attr...)
	return
}
