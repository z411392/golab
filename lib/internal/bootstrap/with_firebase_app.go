package bootstrap

import (
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var (
	firebaseApp *firebase.App
)

func WithFirebaseApp() (app *firebase.App, err error) {
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
		return
	}
	firebaseApp = app
	return
}
