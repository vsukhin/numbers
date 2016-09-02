package logger

import (
	"log"
)

// Logger is an interface of logger
type Logger interface {
	Println(args ...interface{})
	Fatalf(query string, args ...interface{})
}

// DefaultLogger is a default logger implementation
type DefaultLogger struct {
}

var (
	// Log provider Logger
	Log Logger = new(DefaultLogger)
)

// Println prints data to logger
func (logger *DefaultLogger) Println(args ...interface{}) {
	log.Println(args...)
}

// Fatalf prints and stops the execution
func (logger *DefaultLogger) Fatalf(query string, args ...interface{}) {
	log.Fatalf(query, args...)
}
