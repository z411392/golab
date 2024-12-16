package integration_test

import (
	"os"
	"testing"

	"github.com/z411392/golab/container"
)

func TestMain(m *testing.M) {
	container.Init()
	code := m.Run()
	container.Release()
	os.Exit(code)
}
