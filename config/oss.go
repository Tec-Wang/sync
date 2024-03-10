package config

type OssConfig struct {
	AK         string `yaml:"ak"`
	SK         string `yaml:"sk"`
	Endpoint   string `yaml:"endpoint"`
	BucketName string `yaml:"bucket"`
}
