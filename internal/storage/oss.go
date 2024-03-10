package storage

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
)

type OSS struct {
	cli    *oss.Client
	bucket *oss.Bucket
}

func NewOSSClient(ak, sk, endpoint, bucketName string) *OSS {
	client, err := oss.New(endpoint, ak, sk)
	if err != nil {
		panic(err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		panic(err)
	}

	return &OSS{
		cli:    client,
		bucket: bucket,
	}
}

type FileInfo struct {
	Size         int64
	Md5          string
	Key          string
	Owner        oss.Owner
	RestoreInfo  string
	LastModified int64
}

// 获取指定目录下的全部文件的详细信息
func (o *OSS) ListByPrefix(prefix string) ([]*FileInfo, error) {
	objects, err := o.bucket.ListObjectsV2(oss.Prefix(prefix))
	if err != nil {
		return nil, err
	}

	files := make([]*FileInfo, 0, len(objects.Objects))
	for _, object := range objects.Objects {
		files = append(files, &FileInfo{
			Key:          object.Key,
			Size:         object.Size,
			Md5:          object.ETag,
			Owner:        object.Owner,
			RestoreInfo:  object.RestoreInfo,
			LastModified: object.LastModified.Unix(),
		})
	}
	return files, nil
}

// 存储文件到OSS
func (o *OSS) Put(baseDir, filePath string) error {
	objectKey, err := o.objectKey(baseDir, filePath)
	if err != nil {
		return err
	}

	return o.bucket.PutObjectFromFile(objectKey, filePath)
}

// 下载文件到本地
func (o *OSS) Get(baseDir, filePath string) error {
	fileKey, err := o.objectKey(baseDir, filePath)
	if err != nil {
		return err
	}

	return o.bucket.GetObjectToFile(fileKey, filePath)
}

func (o *OSS) objectKey(baseDir, filePath string) (string, error) {
	if filePath == "" {
		return "", errors.New("file path is empty")
	}

	if !strings.HasPrefix(filePath, baseDir) {
		return "", errors.New("file path is not in base dir")
	}

	return filePath[len(baseDir):], nil
}

// 判断文件夹是否存在
func (o *OSS) Exists(path string) (bool, error) {
	return o.bucket.IsObjectExist(path)
}

func (o *OSS) CreateDirIfNotExist(path string) error {
	exists, err := o.Exists(path)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// OSS 创建文件夹需要以 / 结尾
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return o.bucket.PutObject(path, strings.NewReader(""))
}
