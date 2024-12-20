package integreation_test

import (
	"os"
	"testing"

	"gitlab.com/z411392/tt-agent/internal/bootstrap"
)

func TestMain(m *testing.M) {
	code := m.Run()
	bootstrap.Cancel()
	os.Exit(code)
}
