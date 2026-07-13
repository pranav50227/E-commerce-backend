package logger

import (
	"log"
	"os"
)

var infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

// Info logs messages to standard output
func Info(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

// Error logs error messages to standard error
func Error(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}
