package logger

import (
	log "github.com/Sirupsen/logrus"
	"io"
)

func NewLogger(format log.Formatter,level log.Level, writer io.Writer) {
	log.SetFormatter(format)
	log.SetOutput(writer)
	log.SetLevel(level)
}

func LogS3Infof(message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service":  "AWS",
		"Resource": "S3",
	}).Infof(message, a)
}

func LogS3Debugf(message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service":  "AWS",
		"Resource": "S3",
	}).Debugf(message, a)
}

func LogS3Warningf(message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service":  "AWS",
		"Resource": "S3",
	}).Warnf(message, a)
}

func LogS3Errorf(message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service":  "AWS",
		"Resource": "S3",
	}).Errorf(message, a)
}
