package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(format string, v ...any)

	Err() *log.Logger
	Output(calldepth int, s string)
}

type logger struct {
	err  *log.Logger
	info *log.Logger
}

func New() Logger {
	return &logger{
		err:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		info: log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}
}

func (l *logger) Info(format string, v ...any) {
	l.info.Printf(format, v...)
}

func (l *logger) Err() *log.Logger {
	return l.err
}

func (l *logger) Output(calldepth int, s string) {
	l.err.Output(calldepth, s)
}
