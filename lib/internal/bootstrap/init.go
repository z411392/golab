package bootstrap

import (
	"log"

	"context"

	"go.uber.org/dig"
)

var (
	Container *dig.Container
	Context   context.Context
	Cancel    context.CancelFunc
)

func init() {
	Container = dig.New()
	Context, Cancel = context.WithCancel(context.Background())
	if err := Container.Provide(WithFirebaseApp); err != nil {
		log.Fatalf("%s", err.Error())
	}
	if err := Container.Provide(WithFirestore); err != nil {
		log.Fatalf("%s", err.Error())
	}
	if err := Container.Provide(WithStorage); err != nil {
		log.Fatalf("%s", err.Error())
	}
	if err := Container.Provide(WithBucket); err != nil {
		log.Fatalf("%s", err.Error())
	}
	go (func() {
		<-Context.Done()
		firestoreClient.Close()
	})()
}
