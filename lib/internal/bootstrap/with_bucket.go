package bootstrap

import (
	cloudStorage "cloud.google.com/go/storage"
	"firebase.google.com/go/v4/storage"
)

var (
	bucket *cloudStorage.BucketHandle
)

func WithBucket(client *storage.Client) (handle *cloudStorage.BucketHandle, err error) {
	if bucket != nil {
		handle = bucket
		return
	}
	handle, err = client.DefaultBucket()
	if err != nil {
		return
	}
	bucket = handle
	return
}
