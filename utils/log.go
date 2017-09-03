package utils

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	prefix  string
	logfile *os.File
}

func NewLogger(prefix string) *Logger {
	wd, _ := os.Getwd()
	filename := prefix + ".log"

	file, err := os.OpenFile(wd+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Can't open log file for " + filename + " " + err.Error())
	}
	return &Logger{
		prefix:  prefix,
		logfile: file,
	}
}

func (log *Logger) Info(msg ...interface{}) {
	at := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%s - %s - ", at, log.prefix)
	for _, str := range msg {
		line += fmt.Sprintf("%v ", str)
	}
	log.logfile.WriteString(line + "\n")
}

func (log *Logger) Write(p []byte) (n int, err error) {
	n = len(p)
	err = nil
	log.Printf(string(p))
	return
}

func (log *Logger) Printf(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args)
	log.logfile.WriteString(line)
	fmt.Print(line)
}
