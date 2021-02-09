// PACKAGE LOGGER CONFIGURES LOGGER
package logger

import (
	"log"
	"os"
)

// LOGGERS FOR DIFFERENT LOGGING MODE

var DebugLogger *log.Logger //  DEBUG MODE
var ErrorLogger *log.Logger // ERROR MODE
var FatalLogger *log.Logger // FATAL MODE

func init() {
	logFile, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	DebugLogger = log.New(logFile, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(logFile, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLogger = log.New(logFile, "Fatal: ", log.Ldate|log.Ltime|log.Lshortfile)
}
