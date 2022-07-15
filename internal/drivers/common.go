package drivers

import (
	"backup/pkg/archive"
	"context"
)

type (
	Driver interface {
		Backup(ctx context.Context, dir string, archiveType archive.Type) error
	}
)
