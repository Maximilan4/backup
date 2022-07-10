package main

import (
	"backup/internal/backup"
	"context"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	err := backup.Run(ctx)
	if err != nil {
		logrus.Error(err)
	}
}