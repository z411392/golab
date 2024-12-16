package unit_test

import (
	"testing"

	resources_test "github.com/z411392/golab/tests/resources"
	"github.com/z411392/golab/utils"
)

func Test_要能夠推論出物件相應的儲存路徑(t *testing.T) {
	t.SkipNow()
	path, err := utils.PathFor(resources_test.Jpeg_f42227d175f04f59bcb193a25f719cdc, 0)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if path != "jpg/f42227d175f04f59bcb193a25f719cdc" {
		t.FailNow()
	}
}
