package integreation_test

import (
	"testing"

	cloudstorage "cloud.google.com/go/storage"
	"firebase.google.com/go/v4/storage"
	"gitlab.com/z411392/tt-agent/internal/bootstrap"
	"gitlab.com/z411392/tt-agent/internal/utils"
)

func Test_要能初始化Storage(t *testing.T) {
	// t.SkipNow()
	bootstrap.Container.Invoke(func(client *storage.Client) (err error) {
		if client == nil {
			t.FailNow()
		}
		return
	})
}

func Test_要能初始化Bucket(t *testing.T) {
	// t.SkipNow()
	bootstrap.Container.Invoke(func(bucket *cloudstorage.BucketHandle) (err error) {
		if bucket == nil {
			t.FailNow()
		}
		return
	})
}

func Test_要能夠上傳文件(t *testing.T) {
	// t.SkipNow()
	if _, err := utils.UploadFromString(bootstrap.Context, "test", "hello, world", nil); err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func Test_要能夠判斷檔案是否存在(t *testing.T) {
	// t.SkipNow()
	exists, err := utils.ExistObject(bootstrap.Context, "test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if !exists {
		t.FailNow()
	}
}

func Test_要能夠下載檔案內容(t *testing.T) {
	// t.SkipNow()
	content, err := utils.DownloadAsString(bootstrap.Context, "test")
	if err != err {
		t.Fatalf("%s", err.Error())
	}
	if content == "" {
		t.FailNow()
	}
}

func Test_要能夠取得檔案網址(t *testing.T) {
	// t.SkipNow()
	objectUrl, err := utils.GetObjectUrl(bootstrap.Context, "test", 0)
	if err != err {
		t.Fatalf("%s", err.Error())
	}
	if objectUrl == "" {
		t.FailNow()
	}
	t.Logf("%s", utils.GetObjectCdnUrl("test"))
}

func Test_要能夠刪除檔案(t *testing.T) {
	// t.SkipNow()
	if err := utils.DeleteObject(bootstrap.Context, "test"); err != nil {
		t.Fatalf("%s", err.Error())
	}
}
