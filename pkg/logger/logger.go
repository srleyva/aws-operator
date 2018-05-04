package logger

import (
	log "github.com/Sirupsen/logrus"
	"io"
)

func NewLogger(format log.Formatter, level log.Level, writer io.Writer) {
	log.SetFormatter(format)
	log.SetOutput(writer)
	log.SetLevel(level)
}

func LogAWSInfof(resource, message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service":  "AWS",
		"Resource": resource,
	}).Infof(message, a)
}

func LogAWSDebugf(resource, message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service":  "AWS",
		"Resource": resource,
	}).Debugf(message, a)
}

func LogAWSWarningf(resource, message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service":  "AWS",
		"Resource": resource,
	}).Warnf(message, a)
}

func LogAWSErrorf(resource, message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service":  "AWS",
		"Resource": resource,
	}).Errorf(message, a)
}
