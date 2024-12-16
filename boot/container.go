package boot

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	cloudStorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"go.uber.org/dig"
	"google.golang.org/api/option"
)

var Container *dig.Container

var Init func()

var (
	Context context.Context
	cancel  context.CancelFunc
)

func Release() {
	go cancel()
	<-Context.Done()
	release()
}

func init() {
	Context, cancel = context.WithCancel(context.Background())
	Init = sync.OnceFunc(register)
}

var (
	firebaseApp         *firebase.App
	firestoreClient     *firestore.Client
	storageClient       *storage.Client
	defaultBucketHandle *cloudStorage.BucketHandle
)

func register() {
	defer (func() {
		if err := recover(); err != nil {
			release()
		}
	})()

	Container = dig.New()

	Container.Provide(func() (app *firebase.App, err error) {
		if firebaseApp != nil {
			app = firebaseApp
			return
		}
		config := &firebase.Config{
			StorageBucket: os.Getenv("CLOUD_STORAGE_DEFAULT_BUCKET"),
		}
		opts := option.WithCredentialsJSON([]byte(os.Getenv("SERVICE_ACCOUNT_CREDENTIALS")))
		app, err = firebase.NewApp(Context, config, opts)
		if err != nil {
			panic(err)
		}
		firebaseApp = app
		return
	})

	Container.Provide(func(app *firebase.App) (client *firestore.Client, err error) {
		if firestoreClient != nil {
			client = firestoreClient
			return
		}
		client, err = app.Firestore(Context)
		if err != nil {
			panic(err)
		}
		firestoreClient = client
		return
	})

	Container.Provide(func(app *firebase.App) (client *storage.Client, err error) {
		if storageClient != nil {
			client = storageClient
			return
		}
		client, err = app.Storage(Context)
		if err != nil {
			panic(err)
		}
		storageClient = client
		return
	})

	Container.Provide(func(client *storage.Client) (handle *cloudStorage.BucketHandle, err error) {
		if defaultBucketHandle != nil {
			handle = defaultBucketHandle
			return
		}
		handle, err = client.DefaultBucket()
		if err != nil {
			panic(err)
		}
		defaultBucketHandle = handle
		return
	})
}

func release() {
	if firestoreClient != nil {
		firestoreClient.Close()
	}
}
