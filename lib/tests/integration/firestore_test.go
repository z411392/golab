package integreation_test

import (
	"testing"
	"time"

	"cloud.google.com/go/firestore"

	"gitlab.com/z411392/tt-agent/internal/bootstrap"
	"gitlab.com/z411392/tt-agent/internal/utils"
)

const testCollection = "development_tests"
const testDocument = "test"

func Test_要能初始化Firestore(t *testing.T) {
	// t.SkipNow()
	bootstrap.Container.Invoke(func(client *firestore.Client) (err error) {
		if client == nil {
			t.FailNow()
		}
		return
	})
}

type TestData struct {
	Hello string `firestore:"hello"`
}

func Test_要能夠新增紀錄(t *testing.T) {
	t.SkipNow()
	err := bootstrap.Container.Invoke(func(client *firestore.Client) (err error) {
		documentReference := client.Collection(testCollection).Doc(testDocument)
		data := map[string]string{
			"hello": "world",
		}
		_, err = documentReference.Set(bootstrap.Context, data, firestore.MergeAll)
		return
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func Test_要能夠讀取紀錄(t *testing.T) {
	t.SkipNow()
	err := bootstrap.Container.Invoke(func(client *firestore.Client) (err error) {
		documentReference := client.Collection(testCollection).Doc(testDocument)
		documentSnapshot, err := documentReference.Get(bootstrap.Context)
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

func Test_要能夠刪除紀錄(t *testing.T) {
	t.SkipNow()
	err := bootstrap.Container.Invoke(func(client *firestore.Client) (err error) {
		documentReference := client.Collection(testCollection).Doc(testDocument)
		_, err = documentReference.Delete(bootstrap.Context)
		return
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func Test_要能夠監聽資料庫(t *testing.T) {
	// t.SkipNow()
	changeKindCh := make(chan firestore.DocumentChangeKind)
	documentChangeKinds := []firestore.DocumentChangeKind{firestore.DocumentAdded, firestore.DocumentRemoved, firestore.DocumentModified}
	stop, err := utils.OnDocumentsChanged(bootstrap.Context, testCollection, func(changeKind firestore.DocumentChangeKind, parameters []string, data TestData) (err error) {
		switch changeKind {
		case firestore.DocumentAdded:
			// t.Logf("新增了文件")
		case firestore.DocumentModified:
			// t.Logf("更新了文件")
		case firestore.DocumentRemoved:
			// t.Logf("刪除了文件")
		}
		changeKindCh <- changeKind
		return
	}, documentChangeKinds)
	defer stop()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	err = bootstrap.Container.Invoke(func(client *firestore.Client) (err error) {
		collectionReference := client.Collection(testCollection)
		documentReference := collectionReference.Doc(testDocument)
		documentReference.Create(bootstrap.Context, map[string]string{
			"hello": "world",
		})
		update := []firestore.Update{{Path: "hello", Value: time.Now().String()}}
		documentReference.Update(bootstrap.Context, update)
		documentReference.Delete(bootstrap.Context)
		return
	})
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	timeoutCh := time.After(time.Second * 3)
	changeKinds := make([]firestore.DocumentChangeKind, 0)
LOOP:
	for {
		select {
		case changeKind := <-changeKindCh:
			changeKinds = append(changeKinds, changeKind)
			continue LOOP
		case <-timeoutCh:
			break LOOP
		}
	}
	if len(changeKinds) < 3 {
		t.FailNow()
	}
}
