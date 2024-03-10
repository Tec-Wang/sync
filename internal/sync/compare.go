package sync

import (
	"errors"
	"oss-sync/utils"
)

func (s *Sync) fileListAfterCompareDir(localDir string, remoteDir string) (interface{}, error) {
	exists, err := s.ossStorage.Exists(remoteDir)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("remote dir not exists")
	}

	if exists := utils.Exists(localDir); !exists {
		return nil, errors.New("local dir not exists")
	}

	localDirInfo, err := s.scanner.ScanLocalDir(localDir)
	if err != nil {
		return nil, err
	}

	ossDirInfo, err := s.scanner.ScanOssDir(remoteDir)
	if err != nil {
		return nil, err
	}

	return s.scanner.CompareDir(localDirInfo, ossDirInfo)
}
