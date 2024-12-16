package integreation_test

import (
	"testing"

	cloudstorage "cloud.google.com/go/storage"
	"firebase.google.com/go/v4/storage"
	"github.com/z411392/golab/boot"
	"github.com/z411392/golab/utils"
)

func Test_ShouldInitializeStorage(t *testing.T) {
	t.SkipNow()
	boot.Container.Invoke(func(client *storage.Client) (err error) {
		if client == nil {
			t.FailNow()
		}
		return
	})
}

func Test_ShouldInitializeBucket(t *testing.T) {
	t.SkipNow()
	boot.Container.Invoke(func(bucket *cloudstorage.BucketHandle) (err error) {
		if bucket == nil {
			t.FailNow()
		}
		return
	})
}

func Test_ShouldUploadFile(t *testing.T) {
	t.SkipNow()
	_, err := utils.UploadFromString(boot.Context, "test", "hello, world", nil)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func Test_ShouldCheckIfFileExists(t *testing.T) {
	t.SkipNow()
	exists, err := utils.ExistObject(boot.Context, "test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if !exists {
		t.FailNow()
	}
}

func Test_ShouldDownloadFileContent(t *testing.T) {
	t.SkipNow()
	content, err := utils.DownloadAsString(boot.Context, "test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if content == "" {
		t.FailNow()
	}
}

func Test_ShouldGetFileUrl(t *testing.T) {
	t.SkipNow()
	objectUrl, err := utils.GetObjectUrl(boot.Context, "test", 0)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if objectUrl == "" {
		t.FailNow()
	}
}

func Test_ShouldDeleteFile(t *testing.T) {
	t.SkipNow()
	err := utils.DeleteObject(boot.Context, "test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
