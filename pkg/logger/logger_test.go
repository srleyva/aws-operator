package logger

import (
	"testing"
	"bytes"
	"github.com/Sirupsen/logrus"
	"strings"
)

func TestLogS3Infof(t *testing.T) {
	buffer := bytes.Buffer{}
	NewLogger(&logrus.TextFormatter{}, logrus.InfoLevel, &buffer)
	LogS3Infof("Hello %s", "world")
	if !strings.Contains(buffer.String(), "level=info msg=\"Hello [world]\"") {
		t.Errorf("Wrong response given: %s", buffer.String())
	}
}

func TestLogS3Warningf(t *testing.T) {
	buffer := bytes.Buffer{}
	NewLogger(&logrus.TextFormatter{}, logrus.InfoLevel, &buffer)
	LogS3Warningf("Hello %s", "world")
	if !strings.Contains(buffer.String(), "level=warning msg=\"Hello [world]\"") {
		t.Errorf("Wrong response given: %s", buffer.String())
	}
}

func TestLogS3Debugf(t *testing.T) {
	buffer := bytes.Buffer{}
	NewLogger(&logrus.TextFormatter{}, logrus.DebugLevel, &buffer)
	LogS3Debugf("Hello %s", "world")
	if !strings.Contains(buffer.String(), "level=debug msg=\"Hello [world]\"") {
		t.Errorf("Wrong response given: %s", buffer.String())
	}
}

func TestLogS3Fatalf(t *testing.T) {
	buffer := bytes.Buffer{}
	NewLogger(&logrus.TextFormatter{}, logrus.InfoLevel, &buffer)
	LogS3Errorf("Hello %s", "world")
	if !strings.Contains(buffer.String(), "level=error msg=\"Hello [world]\"") {
		t.Errorf("Wrong response given: %s", buffer.String())
	}
}