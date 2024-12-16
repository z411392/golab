package integration_test

import (
	"testing"

	cloudstorage "cloud.google.com/go/storage"
	"firebase.google.com/go/v4/storage"
	"github.com/z411392/golab/container"
	"github.com/z411392/golab/utils"
)

func Test_要能夠初始化CloudStorage(t *testing.T) {
	t.SkipNow()
	container.Container.Invoke(func(client *storage.Client) (err error) {
		if client == nil {
			t.FailNow()
		}
		return
	})
}

func Test_要能夠操作預設的Bucket(t *testing.T) {
	t.SkipNow()
	container.Container.Invoke(func(bucket *cloudstorage.BucketHandle) (err error) {
		if bucket == nil {
			t.FailNow()
		}
		return
	})
}

func Test_要能夠上傳文件(t *testing.T) {
	t.SkipNow()
	_, err := utils.UploadFromString(container.Context, "test", "hello, world", nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func Test_要能夠檢查物件是否存在(t *testing.T) {
	t.SkipNow()
	exists, err := utils.ExistObject(container.Context, "test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if !exists {
		t.FailNow()
	}
}

func Test_要能夠取得文件內容(t *testing.T) {
	t.SkipNow()
	content, err := utils.DownloadAsString(container.Context, "test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if content == "" {
		t.FailNow()
	}
}

func Test_要能夠取得物件網址(t *testing.T) {
	t.SkipNow()
	objectUrl, err := utils.GetObjectUrl(container.Context, "test", 0)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if objectUrl == "" {
		t.FailNow()
	}
}

func Test_要能夠刪除物件(t *testing.T) {
	t.SkipNow()
	err := utils.DeleteObject(container.Context, "test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
