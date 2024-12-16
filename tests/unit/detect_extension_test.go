package unit_test

import (
	"testing"

	resources_test "github.com/z411392/golab/tests/resources"
	"github.com/z411392/golab/utils"
)

func Test_要能從二進制內容中判斷副檔名(t *testing.T) {
	t.SkipNow()
	extension, err := utils.DetectExtension(resources_test.Jpeg_f42227d175f04f59bcb193a25f719cdc)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if extension != "jpg" {
		t.FailNow()
	}
}
