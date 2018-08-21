package logging

import (
	"github.com/sirupsen/logrus"
)

type customFormatter struct {
	logrus.Formatter
}

// Format creates a log format with UTC time and with the standard fields
func (u customFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()

	e.Data = e.WithFields(standardFields).Data

	return u.Formatter.Format(e)
}
