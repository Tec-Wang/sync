package storage

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"oss-sync/config"
	"testing"
)

var ossConfig = config.GetConfig().GetOSSConfig()

func TestClient(t *testing.T) {
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AK, ossConfig.SK)
	if err != nil {
		t.Error(err)
	}

	lsRes, err := client.ListBuckets()
	if err != nil {
		t.Error(err)
	}

	for _, bucket := range lsRes.Buckets {
		fmt.Println("Buckets:", bucket.Name)
	}
}

var client = NewOSSClient(ossConfig.AK, ossConfig.SK, ossConfig.Endpoint, ossConfig.BucketName)

func TestList(t *testing.T) {
	_, err := client.ListByPrefix("")
	if err != nil {
		t.Error(err)
	}
}

func TestPut(t *testing.T) {
	err := client.Put(".", "test.txt")
	if err != nil {
		t.Error(err)
	}
}
