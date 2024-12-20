package resource_test

import (
	_ "embed"
	"testing"
)

var (
	//go:embed images/f42227d175f04f59bcb193a25f719cdc.jpg
	Jpeg_f42227d175f04f59bcb193a25f719cdc []byte
)

func Test_要能載入測試圖片(t *testing.T) {
	if Jpeg_f42227d175f04f59bcb193a25f719cdc == nil {
		t.FailNow()
	}
}
