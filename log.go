package httpstream

import (
	"fmt"
	"log"
)

const (
	FATAL = 0
	ERROR = 1
	WARN  = 2
	INFO  = 3
	DEBUG = 4
)

var (
	LogLevel      int = ERROR
	logger        *log.Logger
	LogLevelWords map[string]int = map[string]int{"fatal": 0, "error": 1, "warn": 2, "info": 3, "debug": 4, "none": -1}
)

// set the logger to use, normal usage:
//
//		var logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
//
//		func init() {
//			SetLogger(log.New(os.Stdout, "", log.Ltime|log.Lshortfile), "debug")
//		}
func SetLogger(l *log.Logger, logLevel string) {
	logger = l
	LogLevelSet(logLevel)
}

// sets the log level from a string
func LogLevelSet(level string) {
	if lvl, ok := LogLevelWords[level]; ok {
		LogLevel = lvl
	}
}

// Log at debug level
func Debug(v ...interface{}) {
	if logger != nil && LogLevel >= 4 {
		logger.Output(2, fmt.Sprintln(v...))
	}
}

func Debugf(format string, v ...interface{}) {
	if LogLevel >= 4 {
		DoLog(3, fmt.Sprintf(format, v...), logger)
	}
}

// Log to logger if setup
//    Log(ERROR, "message")
func Log(logLvl int, v ...interface{}) {
	if LogLevel >= logLvl {
		DoLog(3, fmt.Sprintln(v...), logger)
	}
}

// Log to logger if setup
//    Logf(ERROR, "message %d", 20)
func Logf(logLvl int, format string, v ...interface{}) {
	if LogLevel >= logLvl {
		DoLog(3, fmt.Sprintf(format, v...), logger)
	}
}

func DoLog(depth int, msg string, lgr *log.Logger) {
	if lgr != nil {
		lgr.Output(depth, msg)
	}
}
