package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// 生成MD5
func Md5File(path string) string {
	// 计算文件
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := md5.New()
	io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil))
}

// 创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetFileName(path string) string {
	return filepath.Base(path)
}
