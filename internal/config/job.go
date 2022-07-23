package config

import (
	"backup/internal/drivers"
	"backup/pkg/archive"
	"time"
)

type (
	JobCfg struct {
		Title         string        `yaml:"name" mapstructure:"name"`
		Directory     string        `yaml:"target" mapstructure:"target"`
		DriverInfo    string        `yaml:"driver" mapstructure:"driver"`
		CronSchedule  string        `yaml:"schedule" mapstructure:"schedule"`
		Archive       string        `yaml:"archive" mapstructure:"archive"`
		CancelTimeout time.Duration `yaml:"timeout" mapstructure:"timeout"`
	}
)

func (t *JobCfg) Name() string {
	return t.Title
}

func (t *JobCfg) Target() string {
	return t.Directory
}

func (t *JobCfg) Driver() (drivers.DriverInfo, error) {
	return drivers.NewDriverInfo(t.DriverInfo)
}

func (t *JobCfg) Schedule() string {
	return t.CronSchedule
}

func (t *JobCfg) ArchiveType() archive.Type {
	return archive.Type(t.Archive)
}

func (t *JobCfg) Timeout() time.Duration {
	return t.CancelTimeout
}
