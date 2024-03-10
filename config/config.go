package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const configPath = "/home/neverstop/GolandProjects/oss-sync/config.yaml"

var c = newConfig()

func GetConfig() Config {
	return c
}

func newConfig() Config {
	// 解析Yaml
	dataBytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config config
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		panic(err)
	}

	return &config
}

type Config interface {
	GetOSSConfig() OssConfig
	GetSyncConfig() SyncConfig
}

type config struct {
	OssConfig  OssConfig  `yaml:"oss"`
	SyncConfig SyncConfig `yaml:"sync"`
}

func (c *config) GetOSSConfig() OssConfig {
	return c.OssConfig
}

func (c *config) GetSyncConfig() SyncConfig {
	return c.SyncConfig
}
