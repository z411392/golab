package unit_test

import (
	"testing"

	resources_test "github.com/z411392/golab/tests/resources"
	"github.com/z411392/golab/utils"
)

func Test_要能從二進制內容中判斷MimeType(t *testing.T) {
	t.SkipNow()
	mime, err := utils.DetectMimeType(resources_test.Jpeg_f42227d175f04f59bcb193a25f719cdc)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if mime != "image/jpeg" {
		t.FailNow()
	}
}
