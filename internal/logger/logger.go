package logger

import (
	"fmt"
	"github.com/TwiN/go-color"
	"log"
	"os"
)

type Logger struct {
	l *log.Logger
}

var Log = &Logger{
	l: log.New(os.Stderr, "", 0),
}

func (logger *Logger) Debug(format string, a ...any) {
	text := fmt.Sprintf(format, a)
	logger.l.Printf(color.InCyan(text))
}
