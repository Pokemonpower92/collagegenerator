package service

import (
	"fmt"
	"log"
	"os"
)

// ServiceLogger is a standard logger
// with a tag of the form: [SERVICE] [{servicename}]
type ServiceLogger struct {
	logger *log.Logger
	tag    string
}

// Instantiate a new ServiceLogger that will tag a
// logger like so: [SERVICE] [{service}]
func NewServiceLogger(service string) *ServiceLogger {
	return &ServiceLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
		tag:    fmt.Sprintf("[SERVICE] [%s] ", service),
	}
}

func (l *ServiceLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(l.tag+format, v...)
}

func (l *ServiceLogger) Println(v ...interface{}) {
	args := append([]interface{}{l.tag}, v...)
	l.logger.Println(args...)
}

func (l *ServiceLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(l.tag+format, v...)
}
