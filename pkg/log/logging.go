package log

import "log"

var (
	currentLogLevel uint = 0
)

// SetLogLevel :
func SetLogLevel(level uint) {
	currentLogLevel = level
}

// Logf :
func Logf(level uint, a string, args ...interface{}) {
	if level <= currentLogLevel {
		log.Printf(a, args...)
	}

}

// Log :
func Log(level uint, args ...interface{}) {
	if level <= currentLogLevel {
		log.Print(args...)
	}
}

// Logln :
func Logln(level uint, args ...interface{}) {
	if level <= currentLogLevel {
		log.Println(args...)
	}
}
