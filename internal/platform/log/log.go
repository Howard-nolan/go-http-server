package log

import (
	"fmt"
	"log"
)

type Logger struct{}

func New() *Logger { return &Logger{} }

func (l *Logger) Infof(format string, v ...any)  { log.Printf("INFO  "+format, v...) }
func (l *Logger) Errorf(format string, v ...any) { log.Printf("ERROR "+format, v...) }

// Optional convenience
func (l *Logger) Println(v ...any)          { log.Println(v...) }
func (l *Logger) Printf(f string, v ...any) { log.Printf(f, v...) }

// Example usage: logger.Infof("starting on :%d", port)
func Sprintf(f string, v ...any) string { return fmt.Sprintf(f, v...) }
