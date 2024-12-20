package utils_test

import (
	"testing"

	"gitlab.com/z411392/tt-agent/internal/utils"
	resource_test "gitlab.com/z411392/tt-agent/tests/resource"
)

func Test_要能取判斷二進制內容的檔名(t *testing.T) {
	// t.SkipNow()
	extension, err := utils.DetectExtension(resource_test.Jpeg_f42227d175f04f59bcb193a25f719cdc)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if extension != "jpg" {
		t.FailNow()
	}
}
