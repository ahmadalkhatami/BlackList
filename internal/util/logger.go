package util

import "log"

func Info(msg string, args ...interface{}) {
	log.Printf("Info: "+msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Printf("Error: "+msg, args...)
}
