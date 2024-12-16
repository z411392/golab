package utils

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/z411392/golab/container"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func parametersFromRefPath(path string) []string {
	pathParts := strings.Split(path, "/")
	parameters := make([]string, 0)
	for i := 1; i < len(pathParts); i += 2 {
		parameters = append(parameters, pathParts[i])
	}
	return parameters
}

type NewQuerySnapshotArrived[T any] struct {
	ChangeKind firestore.DocumentChangeKind
	Parameters []string
	Data       T
}

func WhenNewQuerySnapshotArrived[T any](ctx context.Context, q interface{}, documentChangeKinds ...firestore.DocumentChangeKind) (chan NewQuerySnapshotArrived[T], error) {
	var ch chan NewQuerySnapshotArrived[T]
	err := container.Container.Invoke(func(client *firestore.Client) (err error) {
		if collection, ok := q.(string); ok {
			q = client.Collection(collection).Query
		}
		q, ok := q.(firestore.Query)
		if !ok {
			err = fmt.Errorf("expect firestore.Query to be firestore.Query, but got %T", q)
			return
		}
		ch = make(chan NewQuerySnapshotArrived[T])
		startedAt := time.Now()
		it := q.Snapshots(ctx)
		go (func() {
			for {
				select {
				case <-ctx.Done():
					close(ch)
					return
				default:
					snapshot, err := it.Next()
					if status.Code(err) == codes.DeadlineExceeded {
						break
					}
					if err != nil {
						continue
					}
					if snapshot == nil {
						continue
					}
					for _, change := range snapshot.Changes {
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
							continue
						}
						parameters := parametersFromRefPath(change.Doc.Ref.Path)
						ch <- NewQuerySnapshotArrived[T]{
							ChangeKind: change.Kind,
							Parameters: parameters,
							Data:       data,
						}
					}
				}
			}
		})()
		return
	})
	return ch, err
}
