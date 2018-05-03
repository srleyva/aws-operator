package logger

import (
	log "github.com/Sirupsen/logrus"
	"io"
)

func NewLogger(level log.Level, writer io.Writer) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(writer)
	log.SetLevel(level)
}

func LogS3Info(message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service": "AWS",
		"Resource": "S3",
	}).Infof(message, a)
}

func LogS3Debug(message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service": "AWS",
		"Resource": "S3",
	}).Debugf(message, a)
}

func LogS3Warning(message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service": "AWS",
		"Resource": "S3",
	}).Warnf(message, a)
}

func LogS3Error(message string, a ...interface{}) {
	log.WithFields(log.Fields{
		"Service": "AWS",
		"Resource": "S3",
	}).Errorf(message, a)
}
