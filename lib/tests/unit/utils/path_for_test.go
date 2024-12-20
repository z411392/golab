package utils_test

import (
	"testing"

	"gitlab.com/z411392/tt-agent/internal/utils"
	resource_test "gitlab.com/z411392/tt-agent/tests/resource"
)

func Test_要能夠計算儲存的路徑(t *testing.T) {
	// t.SkipNow()
	path, err := utils.PathFor(resource_test.Jpeg_f42227d175f04f59bcb193a25f719cdc, 0)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if path != "jpg/f42227d175f04f59bcb193a25f719cdc" {
		t.FailNow()
	}
}
