package drivers


type (
	S3Driver struct {
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Bucket string `mapstructure:"bucket"`
		Url string `mapstructure:"url"`
		Region string `mapstructure:"region"`
	}
)