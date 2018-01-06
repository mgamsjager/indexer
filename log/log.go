package log

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(msg string)
	Infof(format string, msg ...interface{})
	Error(msg ...interface{})
	Fatal(msg ...interface{})
}

type Logging struct {
}

func New() *Logging {
	return &Logging{}
}

func (l *Logging) Info(msg string) {
	fmt.Printf("\n%v", msg)
}

func (l *Logging) Infof(format string, msg ...interface{}) {
	fmt.Printf(format, msg...)
}

func (l *Logging) Error(msg ...interface{}) {
	log.Println(msg...)
}

func (l *Logging) Fatal(msg ...interface{}) {
	log.Fatalln(msg...)
}
