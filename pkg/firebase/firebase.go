package firebase

import (
	"context"
	"encoding/json"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func init() {
	credentials, err := json.Marshal(map[string]string{
		"type": "service_account",
	})
	if err != nil {
		panic(fmt.Sprintf("無法產生 firebase credentials: %s", err.Error()))
	}
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(credentials))
	if err != nil {
		panic(fmt.Sprintf("無法正確產生 firebase app instance: %s", err.Error()))
	}
}
