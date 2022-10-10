package main

import (
	"backup/internal/backup"
	"context"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	sigCtx, done := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	defer done()

	err := backup.Run(sigCtx)
	if err != nil {
		logrus.Error(err)
	}
}
