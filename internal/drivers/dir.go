package drivers

const (
	DirectoryType = "dir"
)

type (
	DirectoryDriver struct {
		Path string `mapstructure:"path"`
	}
)