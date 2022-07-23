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
	c.logger.WithField("info", keysAndValues).Info(msg)
}

func (c *CronAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	c.logger.WithField("info", keysAndValues).WithField("msg", msg).Error(err)
}
