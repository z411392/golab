package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/h2non/filetype"
	"gitlab.com/z411392/tt-agent/internal/bootstrap"
)

func DetectMimeType(buffer []byte) (string, error) {
	kind, err := filetype.Match(buffer)
	if err != nil {
		return "", err
	}
	return kind.MIME.Value, nil
}

func DetectExtension(buffer []byte) (string, error) {
	kind, err := filetype.Match(buffer)
	if err != nil {
		return "", err
	}
	return kind.Extension, nil
}

func existObject(ctx context.Context, path string) (exists bool, err error) {
	bootstrap.Container.Invoke(func(bucket *storage.BucketHandle) {
		prefix := Prefix()
		object := bucket.Object(fmt.Sprintf("%s/%s", prefix, path))
		_, err = object.Attrs(ctx)
		switch {
		case err == storage.ErrObjectNotExist:
			err = nil
			return
		case err == nil:
			exists = true
			return
		default:
			return
		}
	})
	return
}

func ExistObject(ctx context.Context, path string) (exists bool, err error) {
	existsCh := make(chan bool)
	errCh := make(chan error)
	go (func() {
		exists, err = existObject(ctx, path)
		if err == nil {
			existsCh <- exists
		} else {
			errCh <- err
		}
		close(existsCh)
		close(errCh)
	})()
	select {
	case exists = <-existsCh:
		break
	case err = <-errCh:
		break
	}
	return
}

func putObject(ctx context.Context, path string, buffer []byte, metadata map[string]string) (bytes int, err error) {
	bootstrap.Container.Invoke(func(bucket *storage.BucketHandle) {
		prefix := Prefix()
		object := bucket.Object(fmt.Sprintf("%s/%s", prefix, path))
		mime, _ := DetectMimeType(buffer)
		writer := object.NewWriter(ctx)
		defer writer.Close()
		if mime != "" {
			writer.ContentType = mime
		}
		if metadata != nil {
			writer.Metadata = metadata
		}
		bytes, err = writer.Write(buffer)
		if err != nil {
			return
		}
	})
	return
}

func PutObject(ctx context.Context, path string, buffer []byte, metadata map[string]string) (bytes int, err error) {
	bytesCh := make(chan int)
	errCh := make(chan error)
	go (func() {
		bytes, err := putObject(ctx, path, buffer, metadata)
		if err == nil {
			bytesCh <- bytes
		} else {
			errCh <- err
		}
		close(errCh)
		close(bytesCh)
	})()
	select {
	case err = <-errCh:
		break
	case bytes = <-bytesCh:
		break
	}
	return
}

func UploadFromString(ctx context.Context, path string, content string, metadata map[string]string) (bytes int, err error) {
	bytes, err = PutObject(ctx, path, []byte(content), metadata)
	return
}

func Prefix() (projectId string) {
	projectId = os.Getenv("FIREBASE_PROJECT_ID")
	env := "production"
	if OnDevelopment() {
		env = "development"
	}
	return fmt.Sprintf("%s/%s", projectId, env)
}

func PathFor(buffer []byte, bufferSize int) (path string, err error) {
	if bufferSize == 0 {
		bufferSize = 4096
	}
	extension, err := DetectExtension(buffer)
	if err != nil {
		return
	}
	hash := md5.New()
	for i := 0; i < len(buffer); i += bufferSize {
		end := i + bufferSize
		if end > len(buffer) {
			end = len(buffer)
		}
		hash.Write(buffer[i:end])
	}
	digest := hex.EncodeToString(hash.Sum(nil))
	path = fmt.Sprintf("%s/%s", extension, digest)
	return
}

func getObjectUrl(_ context.Context, path string, expiry time.Duration) (objectUrl string, err error) {
	if expiry == 0 {
		expiry = time.Hour
	}
	bootstrap.Container.Invoke(func(bucket *storage.BucketHandle) {
		prefix := Prefix()
		opts := &storage.SignedURLOptions{
			Method:  "GET",
			Expires: time.Now().Add(expiry),
		}
		objectUrl, err = bucket.SignedURL(fmt.Sprintf("%s/%s", prefix, path), opts)
	})
	return
}

func GetObjectUrl(ctx context.Context, path string, expiry time.Duration) (objectUrl string, err error) {
	errCh := make(chan error)
	objectUrlCh := make(chan string)
	go (func() {
		objectUrl, err := getObjectUrl(ctx, path, expiry)
		if err == nil {
			objectUrlCh <- objectUrl
		} else {
			errCh <- err
		}
		close(errCh)
		close(objectUrlCh)
	})()
	select {
	case err = <-errCh:
		break
	case objectUrl = <-objectUrlCh:
		break
	}
	return
}

func GetObjectCdnUrl(path string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("CDN_BASE_URL"), path)
}

func downloadAsString(ctx context.Context, path string) (content string, err error) {
	bootstrap.Container.Invoke(func(bucket *storage.BucketHandle) {
		prefix := Prefix()
		object := bucket.Object(fmt.Sprintf("%s/%s", prefix, path))
		reader, err := object.NewReader(ctx)
		if err != nil {
			return
		}
		defer reader.Close()
		var buffer bytes.Buffer
		if _, err := io.Copy(&buffer, reader); err != nil {
			return
		}
		content = buffer.String()
	})
	return
}

func DownloadAsString(ctx context.Context, path string) (content string, err error) {
	contentCh := make(chan string)
	errCh := make(chan error)
	go (func() {
		content, err := downloadAsString(ctx, path)
		if err == nil {
			contentCh <- content
		} else {
			errCh <- err
		}
		close(errCh)
		close(contentCh)
	})()
	select {
	case err = <-errCh:
		break
	case content = <-contentCh:
		break
	}
	return
}

func deleteObject(ctx context.Context, path string) (err error) {
	bootstrap.Container.Invoke(func(bucket *storage.BucketHandle) {
		prefix := Prefix()
		object := bucket.Object(fmt.Sprintf("%s/%s", prefix, path))
		err = object.Delete(ctx)
	})
	return
}

func DeleteObject(ctx context.Context, path string) (err error) {
	errCh := make(chan error)
	go (func() {
		err := deleteObject(ctx, path)
		errCh <- err
		close(errCh)
	})()
	err = <-errCh
	return
}
