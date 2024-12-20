package utils_test

import (
	"os"
	"testing"

	"gitlab.com/z411392/tt-agent/internal/utils"
)

func Test_要能夠判斷是否在開發環境(t *testing.T) {
	// t.SkipNow()
	os.Setenv("ENV", "")
	if utils.OnDevelopment() {
		t.FailNow()
	}
	os.Setenv("ENV", "development")
	if !utils.OnDevelopment() {
		t.FailNow()
	}
}
