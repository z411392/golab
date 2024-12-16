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
	"github.com/z411392/golab/container"
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

func ExistObject(ctx context.Context, path string) (exists bool, err error) {
	container.Container.Invoke(func(bucket *storage.BucketHandle) {
		prefix := Prefix()
		object := bucket.Object(fmt.Sprintf("%s/%s", prefix, path))
		_, err = object.Attrs(ctx)
		switch {
		case err == storage.ErrObjectNotExist:
			err = nil
		case err == nil:
			exists = true
		default:
			return
		}
	})
	return
}

func PutObject(ctx context.Context, path string, buffer []byte, metadata map[string]string) (n int, err error) {
	container.Container.Invoke(func(bucket *storage.BucketHandle) {
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
		n, err = writer.Write(buffer)
		if err != nil {
			return
		}
	})
	return
}

func UploadFromString(ctx context.Context, path string, content string, metadata map[string]string) (n int, err error) {
	return PutObject(ctx, path, []byte(content), metadata)
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

func GetObjectUrl(_ context.Context, path string, expiry time.Duration) (objectUrl string, err error) {
	if expiry == 0 {
		expiry = time.Hour
	}
	container.Container.Invoke(func(bucket *storage.BucketHandle) {
		prefix := Prefix()
		opts := &storage.SignedURLOptions{
			Method:  "GET",
			Expires: time.Now().Add(expiry),
		}
		objectUrl, err = bucket.SignedURL(fmt.Sprintf("%s/%s", prefix, path), opts)
	})
	return
}

func GetObjectCdnUrl(path string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("CDN_BASE_URL"), path)
}

func DownloadAsString(ctx context.Context, path string) (content string, err error) {
	container.Container.Invoke(func(bucket *storage.BucketHandle) {
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

func DeleteObject(ctx context.Context, path string) (err error) {
	container.Container.Invoke(func(bucket *storage.BucketHandle) {
		prefix := Prefix()
		object := bucket.Object(fmt.Sprintf("%s/%s", prefix, path))
		err = object.Delete(ctx)
	})
	return
}
