package utils_test

import (
	"testing"

	"gitlab.com/z411392/tt-agent/internal/utils"
	resource_test "gitlab.com/z411392/tt-agent/tests/resource"
)

func Test_要能取判斷二進制內容的類型(t *testing.T) {
	// t.SkipNow()
	mime, err := utils.DetectMimeType(resource_test.Jpeg_f42227d175f04f59bcb193a25f719cdc)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	if mime != "image/jpeg" {
		t.FailNow()
	}
}
