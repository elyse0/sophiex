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
	var text string
	if len(a) == 0 {
		text = fmt.Sprintf(format)
	} else {
		text = fmt.Sprintf(format, a)
	}
	logger.l.Printf(color.InCyan(text))
}
