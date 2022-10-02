package logging

import (
	"github.com/sirupsen/logrus"
)

type CronAdapter struct {
	logger *logrus.Logger
}

func NewCronAdapter() *CronAdapter {
	return &CronAdapter{logger: logrus.StandardLogger()}
}

func (c *CronAdapter) Info(msg string, keysAndValues ...interface{}) {
	c.getEntry(keysAndValues).Info(msg)
}

func (c *CronAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	c.getEntry(keysAndValues).WithField("msg", msg).Error(err)
}

func (c *CronAdapter) getEntry(keysAndValues []any) *logrus.Entry {
	entry := logrus.NewEntry(c.logger)
	for pos := 0; pos < len(keysAndValues); pos = pos + 2 {
		if pos%2 != 0 || pos == len(keysAndValues) {
			continue
		}
		entry = entry.WithField(keysAndValues[pos].(string), keysAndValues[pos+1])
	}

	return entry
}
