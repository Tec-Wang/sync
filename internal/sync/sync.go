package sync

import (
	"oss-sync/config"
	"oss-sync/internal/storage"
	"oss-sync/internal/xray"
	"oss-sync/utils"
)

type Sync struct {
	config     config.SyncConfig
	ossStorage *storage.OSS
	scanner    *xray.Scanner
}

func NewSync() *Sync {
	c := config.GetConfig()
	syncConfig := c.GetSyncConfig()
	ossConfig := c.GetOSSConfig()
	client := storage.NewOSSClient(ossConfig.AK, ossConfig.SK, ossConfig.Endpoint, ossConfig.BucketName)
	return &Sync{
		config:     syncConfig,
		scanner:    xray.NewScanner(client),
		ossStorage: client,
	}
}

func (s *Sync) Run() {
	// 校验配置，创建目录
	err := s.initSyncDirs()
	if err != nil {
		panic(err)
	}

	// 同步目录
	for _, syncDir := range s.config.SyncDirs {
		go s.syncDir(syncDir)
	}

	// 同步文件
	for _, syncFile := range s.config.SyncFiles {
		go s.syncFile(syncFile)
	}
}

// 校验远程的目录是否存在，以及本地目录是否存在
// 如果本地目录不存在，则创建
// 如果远程目录不存在，则记录日志，并且报错
func (s *Sync) initSyncDirs() error {
	dirs := s.config.SyncDirs
	for _, dir := range dirs {
		if err := s.initDir(dir); err != nil {
			return err
		}
	}

	return nil
}

func (s *Sync) initRemoteDir(dir string) error {
	return s.ossStorage.CreateDirIfNotExist(dir)
}

func (s *Sync) initDir(dir config.SyncDir) error {
	err := s.initLocalDir(dir.LocalPath)
	if err != nil {
		return err
	}

	return s.initRemoteDir(dir.RemotePath)
}

func (s *Sync) syncDir(dir config.SyncDir) {
	// 生成最终文件列表
	fileList, err := s.fileListAfterCompareDir(dir.LocalPath, dir.RemotePath)
	if err != nil {
		panic(err)
	}

	// 获取所有的本地文件和远端文件，进行文件比较。
	// 维持一个最终的文件列表和对应的索引位置
	// 1. 文件在本地，在远端，且文件内容一致，则不做任何事
	// 2. 文件在本地，在远端，且文件内容不一致，则更新
	// 2.1 本地文件最终更新时间小于远端文件，使用远端文件覆盖本地文件
	// 2.2 本地文件最终更新时间大于远端文件，使用本地文件覆盖远端文件
	// 3. 文件只存在于本地，则上传到远端
	// 4. 文件只存在于远端，则下载到本地
}

func (s *Sync) initLocalDir(dir string) error {
	if exists := utils.Exists(dir); exists {
		return nil
	}

	return utils.CreateDir(dir)
}

func (s *Sync) syncFile(file config.SyncFile) {

}
