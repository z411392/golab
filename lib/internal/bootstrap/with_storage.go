package bootstrap

import (
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
)

var (
	storageClient *storage.Client
)

func WithStorage(app *firebase.App) (client *storage.Client, err error) {
	if storageClient != nil {
		client = storageClient
		return
	}
	client, err = app.Storage(Context)
	if err != nil {
		return
	}
	storageClient = client
	return
}
