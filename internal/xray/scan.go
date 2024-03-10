package xray

import (
	"errors"
	"os"
	"oss-sync/internal/entity"
	"oss-sync/internal/storage"
	"oss-sync/utils"
)

type Scanner struct {
	storage *storage.OSS
}

func NewScanner(storage *storage.OSS) *Scanner {
	return &Scanner{
		storage: storage,
	}
}

// 扫描本地文件，获取文件的MD5值和修改时间
func (s *Scanner) ScanLocalDir(dirPath string) ([]*entity.FileScanInfo, error) {
	res := make([]*entity.FileScanInfo, 0)

	// 判断目标路径是否为文件夹
	if !utils.IsDir(dirPath) {
		return nil, errors.New("dirPath is not a dir")
	}

	// 获取文件夹下所有文件
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, item := range files {
		path := dirPath + "/" + item.Name()

		switch item.Type() {
		case os.ModeDir:
			dirScanRes, err := s.ScanLocalDir(path)
			if err != nil {
				return nil, err
			}
			res = append(res, dirScanRes...)
		default:
			fileInfo, err := s.ScanFile(path)
			if err != nil {
				panic(err)
			}
			res = append(res, fileInfo)
		}

	}

	return res, nil
}

func (*Scanner) ScanFile(filePath string) (*entity.FileScanInfo, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &entity.FileScanInfo{
		Path:     filePath,
		Md5:      utils.Md5File(filePath),
		Size:     stat.Size(),
		UpdateAt: stat.ModTime().Unix(),
	}, nil
}

func (s *Scanner) ScanOssDir(dirPath string) ([]*entity.FileScanInfo, error) {
	objectsInfo, err := s.storage.ListByPrefix(dirPath)
	if err != nil {
		return nil, err
	}

	res := make([]*entity.FileScanInfo, 0, len(objectsInfo))

	for _, item := range objectsInfo {
		res = append(res, &entity.FileScanInfo{
			Path:     item.Key,
			Md5:      item.Md5,
			Size:     item.Size,
			UpdateAt: item.LastModified,
		})
	}

	return res, nil
}

type CompareFileRes struct {
	location  FileLocation
	Path      string
	Operation OperationType
}

type CompareFileList []*CompareFileRes

// CompareDir 比较本地和远端文件列表
func (s *Scanner) CompareDir(localScanInfo []*entity.FileScanInfo, remoteScanInfo []*entity.FileScanInfo) (*CompareFileList, error) {
	localMap := make(map[string]*entity.FileScanInfo)
	remoteMap := make(map[string]*entity.FileScanInfo)

	// todo 这里应该将前缀去掉，有的是文件夹。擦 不想写了。真麻烦
	for _, item := range localScanInfo {
		localMap[item.Path] = item
	}

	for _, item := range remoteScanInfo {
		remoteMap[item.Path] = item
	}

	res := make([]*CompareFileRes, 0, len(localMap))

	// 该复习下设计模式了
	for path, localInfo := range localMap {
		// 获取文件名称
		fileName := utils.GetFileName(path)

		remoteInfo, ok := remoteMap[path]

		if !ok {
			res = append(res, &CompareFileRes{
				location:  Local,
				Path:      path,
				Operation: Delete,
			})
		} else {
			if localInfo.Md5 != remoteInfo.Md5 {
				res = append(res, &CompareFileRes{
					location:  Local,
					Path:      path,
					Operation: Update,
				})
			}
		}
	}

	for path, remoteInfo := range remoteMap {
		_, ok := localMap[path]

		if !ok {
			res = append(res, &CompareFileRes{
				location:  Remote,
				Path:      path,
				Operation: Create,
			})

		}
	}

	return &res, nil

}
