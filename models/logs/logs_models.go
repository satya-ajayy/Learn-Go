package logs

import (
	"learn-go/errors"
	"strings"
)

type LogModel struct {
	ID        string `json:"id" bson:"_id"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
	Level     string `json:"level" bson:"level"`
	Message   string `json:"message" bson:"message"`
	Context   string `json:"context" bson:"context"`
}

func (l *LogModel) Validate() error {
	ve := errors.ValidationErrs()

	if l.Level == "" {
		ve.Add("level", "cannot be empty")
	}
	if l.Message == "" {
		ve.Add("message", "cannot be empty")
	}
	if l.Context == "" {
		ve.Add("context", "cannot be empty")
	}
	l.Level = strings.ToUpper(l.Level)
	return ve.Err()
}
