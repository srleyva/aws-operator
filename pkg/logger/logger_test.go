package logger

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"strings"
	"testing"
)

func TestLogAWSInfof(t *testing.T) {
	buffer := bytes.Buffer{}
	NewLogger(&logrus.TextFormatter{}, logrus.InfoLevel, &buffer)
	LogAWSInfof("S3", "Hello %s", "world")
	if !strings.Contains(buffer.String(), "level=info msg=\"Hello [world]\"") {
		t.Errorf("Wrong response given: %s", buffer.String())
	}
}

func TestLogAWSWarningf(t *testing.T) {
	buffer := bytes.Buffer{}
	NewLogger(&logrus.TextFormatter{}, logrus.InfoLevel, &buffer)
	LogAWSWarningf("S3", "Hello %s", "world")
	if !strings.Contains(buffer.String(), "level=warning msg=\"Hello [world]\"") {
		t.Errorf("Wrong response given: %s", buffer.String())
	}
}

func TestLogAWSDebugf(t *testing.T) {
	buffer := bytes.Buffer{}
	NewLogger(&logrus.TextFormatter{}, logrus.DebugLevel, &buffer)
	LogAWSDebugf("S3", "Hello %s", "world")
	if !strings.Contains(buffer.String(), "level=debug msg=\"Hello [world]\"") {
		t.Errorf("Wrong response given: %s", buffer.String())
	}
}

func TestLogAWSFatalf(t *testing.T) {
	buffer := bytes.Buffer{}
	NewLogger(&logrus.TextFormatter{}, logrus.InfoLevel, &buffer)
	LogAWSErrorf("S3", "Hello %s", "world")
	if !strings.Contains(buffer.String(), "level=error msg=\"Hello [world]\"") {
		t.Errorf("Wrong response given: %s", buffer.String())
	}
}
