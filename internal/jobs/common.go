package jobs

import (
	"backup/internal/drivers"
	"backup/pkg/archive"
	"backup/pkg/filesystem"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	Config interface {
		Name() string
		Target() string
		Driver() (drivers.DriverInfo, error)
		Schedule() string
		ArchiveType() archive.Type
		Timeout() time.Duration
	}
	Job[C Config, D drivers.Driver] struct {
		cfg    C
		driver D
	}
)

func NewJob(cfg Config, driver drivers.Driver) *Job[Config, drivers.Driver] {
	return &Job[Config, drivers.Driver]{cfg, driver}
}

func (j *Job[C, D]) Run() {
	logrus.WithFields(logrus.Fields{
		"taskName":    j.cfg.Name(),
		"dir":         j.cfg.Target(),
		"schedule":    j.cfg.Schedule(),
		"archiveType": j.cfg.ArchiveType(),
	}).Info("running a task")

	path, err := filesystem.NormalizePath(j.cfg.Target())
	if err != nil {
		logrus.Fatal(err)
	}

	ctx, done := context.WithTimeout(context.Background(), 5*time.Minute)
	defer done()
	if err = j.driver.Backup(ctx, path, j.cfg.ArchiveType()); err != nil {
		logrus.Fatal(err)
	}
}
