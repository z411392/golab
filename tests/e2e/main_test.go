package e2e_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/z411392/golab/boot"
)

var handler http.Handler

func TestMain(m *testing.M) {
	boot.Init()
	gin.SetMode(gin.TestMode)
	handler = boot.NewHttpHandler()
	code := m.Run()
	boot.Release()
	os.Exit(code)
}
