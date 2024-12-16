package unit_test

import (
	"os"
	"testing"

	"github.com/z411392/golab/utils"
)

func Test_要能判斷是否在開發環境(t *testing.T) {
	t.SkipNow()
	os.Setenv("ENV", "")
	if utils.OnDevelopment() {
		t.FailNow()
	}
	os.Setenv("ENV", "development")
	if !utils.OnDevelopment() {
		t.FailNow()
	}
}
