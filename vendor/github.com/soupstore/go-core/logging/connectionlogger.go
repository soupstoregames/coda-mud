package logging

import "github.com/sirupsen/logrus"

type ConnectionLogger struct {
	entry *logrus.Entry
}

func BuildConnectionLogger(connectionID string) *ConnectionLogger {
	return &ConnectionLogger{
		entry: logger.WithField("connection-id", connectionID),
	}
}

// Info logs an info level message with standard fields
func (c *ConnectionLogger) Info(msg string) {
	c.entry.Info(msg)
}

// Warn logs an warn level message with standard fields
func (c *ConnectionLogger) Warn(msg string) {
	c.entry.Warn(msg)
}

// Error logs an error level message with standard fields
func (c *ConnectionLogger) Error(msg string) {
	c.entry.Error(msg)
}

// Fatal logs an fatal level message with standard fields
func (c *ConnectionLogger) Fatal(msg string) {
	c.entry.Fatal(msg)
}
