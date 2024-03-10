package config

type SyncDir struct {
	RemotePath string `yaml:"remotePath"`
	LocalPath  string `yaml:"localPath"`
}

type SyncFile struct {
	RemotePath string `yaml:"remotePath"`
	LocalPath  string `yaml:"localPath"`
}

type SyncConfig struct {
	SyncDirs  []SyncDir  `yaml:"syncDirs"`
	SyncFiles []SyncFile `yaml:"syncFiles"`
}
