package logmatic

import (
	"encoding/json"
	"fmt"
	"time"
	log "github.com/Sirupsen/logrus"
	"os"
)

const defaultTimestampFormat = time.RFC3339
var markers = [2]string{"sourcecode", "golang"}

type JSONFormatter struct {
}

func (f *JSONFormatter) Format(entry *log.Entry) ([]byte, error) {

	data := make(log.Fields, len(entry.Data) + 3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/Sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	//prefixFieldClashes(data)

	data["date"] = entry.Time.Format(defaultTimestampFormat)
	data["message"] = entry.Message
	data["level"] = entry.Level.String()
	data["@marker"] = markers
	data["appname"] = os.Args[0]
	h, err := os.Hostname()
	if err == nil {
		data["hostname"] = h
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
