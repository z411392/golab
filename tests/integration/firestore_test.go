package integration_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/z411392/golab/container"
	"github.com/z411392/golab/utils"
)

const testCollection = "development_tests"
const testDocument = "test"

func Test_要能夠初始化Firestore(t *testing.T) {
	t.SkipNow()
	container.Container.Invoke(func(client *firestore.Client) (err error) {
		if client == nil {
			t.FailNow()
		}
		return
	})
}

type TestData struct {
	Hello string `firestore:"hello"`
}

func Test_要能夠寫入資料(t *testing.T) {
	t.SkipNow()
	err := container.Container.Invoke(func(client *firestore.Client) (err error) {
		documentReference := client.Collection(testCollection).Doc(testDocument)
		data := map[string]string{
			"hello": "world",
		}
		_, err = documentReference.Set(container.Context, data, firestore.MergeAll)
		return
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func Test_要能夠取得資料(t *testing.T) {
	t.SkipNow()
	err := container.Container.Invoke(func(client *firestore.Client) (err error) {
		documentReference := client.Collection(testCollection).Doc(testDocument)
		documentSnapshot, err := documentReference.Get(container.Context)
		if err != nil {
			return
		}
		var data TestData
		err = documentSnapshot.DataTo(&data)
		if data.Hello != "world" {
			t.FailNow()
			return
		}
		return
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func Test_要能夠刪除資料(t *testing.T) {
	t.SkipNow()
	err := container.Container.Invoke(func(client *firestore.Client) (err error) {
		documentReference := client.Collection(testCollection).Doc(testDocument)
		_, err = documentReference.Delete(container.Context)
		return
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func Test_要能夠接收到任何改動後的Snapshots(t *testing.T) {
	t.SkipNow()
	ctx, cancel := context.WithTimeout(container.Context, 5*time.Second)
	defer cancel()
	ch, err := utils.WhenNewQuerySnapshotArrived[TestData](ctx, testCollection, firestore.DocumentAdded, firestore.DocumentModified, firestore.DocumentRemoved)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	timeout := make(chan interface{})
	completed := make(chan firestore.DocumentChangeKind)
	go (func() {
		for {
			select {
			case querySnapshotArrived := <-ch:
				switch querySnapshotArrived.ChangeKind {
				case firestore.DocumentAdded:
					completed <- firestore.DocumentAdded
				case firestore.DocumentModified:
					completed <- firestore.DocumentModified
				case firestore.DocumentRemoved:
					completed <- firestore.DocumentRemoved
				}
			case <-ctx.Done():
				timeout <- nil
				close(timeout)
				return
			}
		}
	})()
	container.Container.Invoke(func(client *firestore.Client) {
		collectionReference := client.Collection(testCollection)
		documentReference := collectionReference.Doc(testDocument)
		documentReference.Create(container.Context, map[string]string{
			"hello": "world",
		})
		update := []firestore.Update{{Path: "hello", Value: time.Now().String()}}
		documentReference.Update(container.Context, update)
		documentReference.Delete(container.Context)
	})
	counter := 0
	for {
		select {
		case <-timeout:
			if counter < 3 {
				t.FailNow()
			}
		case <-completed:
			counter += 1
			if counter >= 3 {
				return
			}
		}
	}
}
