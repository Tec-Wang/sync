package xray

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestScanDir(t *testing.T) {
	filesInfo, err := ScanLocalDir("/home/neverstop/code/godemo")
	if err != nil {
		t.Error(err)
	}
	logger := logrus.New()

	// 不掺杂任何日志
	logger.SetFormatter(&logrus.JSONFormatter{})

	for _, fileInfo := range filesInfo {
		logger.WithFields(logrus.Fields{
			"path":     fileInfo.Path,
			"size":     fileInfo.Size,
			"md5":      fileInfo.Md5,
			"updateAt": fileInfo.UpdateAt,
		}).Info("scan file info")
	}
}
