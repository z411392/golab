package firebase_test

import (
	"testing"

	. "gitlab.com/z411392/tt-agent/pkg/firebase"
)

func TestInit(t *testing.T) {
	if FirebaseApp == nil {
		t.Fatal("無法正確產生 firebase app instance")
	}
}
