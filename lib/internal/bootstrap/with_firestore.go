package bootstrap

import (
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

var (
	firestoreClient *firestore.Client
)

func WithFirestore(app *firebase.App) (client *firestore.Client, err error) {
	if firestoreClient != nil {
		client = firestoreClient
		return
	}
	client, err = app.Firestore(Context)
	firestoreClient = client
	return
}
