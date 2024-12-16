package e2e_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/z411392/golab/config"
	"github.com/z411392/golab/container"
)

var handler http.Handler

func TestMain(m *testing.M) {
	container.Init()
	gin.SetMode(gin.TestMode)
	handler = config.NewHttpHandler()
	code := m.Run()
	container.Release()
	os.Exit(code)
}
