package kecpsignal

import (
	stdlog "log"
)

type Logger interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
}

var logger Logger = stdlog.Default()

func SetLogger(newLogger Logger) {
	logger = newLogger
}
