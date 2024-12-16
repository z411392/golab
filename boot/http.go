package boot

import (
	"net/http"

	"github.com/gin-gonic/gin"
	identityAndAccessManagingControllers "github.com/z411392/golab/modules/identity_and_access_managing/presentation/http/controllers"
	identityAndAccessManagingMiddlewares "github.com/z411392/golab/modules/identity_and_access_managing/presentation/http/middlewares"
	systemMaintainingControllers "github.com/z411392/golab/modules/system_maintaining/presentation/http/controllers"
)

func NewHttpHandler() http.Handler {
	router := gin.Default()
	router.Use(identityAndAccessManagingMiddlewares.WithCredentialResolving())
	router.GET("/liveness_check", systemMaintainingControllers.OnCheckingLiveness)
	router.GET("/readiness_check", systemMaintainingControllers.OnCheckingReadiness)
	router.Use(identityAndAccessManagingMiddlewares.EnsureUserSignedIn())
	router.GET("/auth/me", identityAndAccessManagingControllers.OnRetrievingProfile)
	return router
}
