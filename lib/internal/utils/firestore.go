package utils

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"gitlab.com/z411392/tt-agent/internal/bootstrap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DocumentsChangedHandler[T any] func(changeKind firestore.DocumentChangeKind, parameters []string, data T) error

func ParametersOf(path string) []string {
	pathParts := strings.Split(path, "/")
	parameters := make([]string, 0)
	for i := 1; i < len(pathParts); i += 2 {
		parameters = append(parameters, pathParts[i])
	}
	return parameters
}

func onDocumentsChanged[T any](ctx context.Context, q interface{}, onChangedHandler DocumentsChangedHandler[T], documentChangeKinds []firestore.DocumentChangeKind) (stop func(), err error) {
	err = bootstrap.Container.Invoke(func(client *firestore.Client) (err error) {
		if collection, ok := q.(string); ok {
			q = client.Collection(collection).Query
		}
		q, ok := q.(firestore.Query)
		if !ok {
			err = fmt.Errorf("type assertion failed: expected string but got %T", q)
			return
		}
		startedAt := time.Now()
		it := q.Snapshots(ctx)
		stop = it.Stop
		go (func() {
			for {
				querySnap, err := it.Next()
				if status.Code(err) == codes.DeadlineExceeded {
					break
				}
				if err != nil {
					continue
				}
				if querySnap == nil {
					continue
				}
				for _, change := range querySnap.Changes {
					if OnDevelopment() != strings.Contains(change.Doc.Ref.Path, "development_") {
						continue
					}
					if change.Kind == firestore.DocumentModified {
						if slices.Index(documentChangeKinds, firestore.DocumentModified) == -1 {
							continue
						}
						if change.Doc.UpdateTime.Before(startedAt) {
							continue
						}
					}
					if change.Kind == firestore.DocumentAdded {
						if slices.Index(documentChangeKinds, firestore.DocumentAdded) == -1 {
							continue
						}
						if change.Doc.CreateTime.Before(startedAt) {
							continue
						}
					}
					if change.Kind == firestore.DocumentRemoved {
						if slices.Index(documentChangeKinds, firestore.DocumentRemoved) == -1 {
							continue
						}
					}
					var data T
					if err := change.Doc.DataTo(&data); err != nil {
						fmt.Printf("%s", err.Error())
						continue
					}
					parameters := ParametersOf(change.Doc.Ref.Path)
					onChangedHandler(change.Kind, parameters, data)
				}
			}
		})()
		return
	})
	return
}

func OnDocumentsChanged[T any](ctx context.Context, q interface{}, onChangedHandler DocumentsChangedHandler[T], documentChangeKinds []firestore.DocumentChangeKind) (stop context.CancelFunc, err error) {
	errCh := make(chan error)
	stopCh := make(chan func())
	go (func() {
		stop, err := onDocumentsChanged(ctx, q, onChangedHandler, documentChangeKinds)
		if err != nil {
			errCh <- err
			return
		}
		stopCh <- stop
	})()
	select {
	case err = <-errCh:
		break
	case stop = <-stopCh:
		break
	}
	return
}
