package tasks

type (
	Task interface {
		Name() string
		Target() string
		Driver() string
		Schedule() string
		ArchiveType() string
	}
)
