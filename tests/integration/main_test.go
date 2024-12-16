package integreation_test

import (
	"os"
	"testing"

	"github.com/z411392/golab/boot"
)

func TestMain(m *testing.M) {
	boot.Init()
	code := m.Run()
	boot.Release()
	os.Exit(code)
}
