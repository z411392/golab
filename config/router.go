package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// identityAndAccessManagingControllers "github.com/z411392/golab/modules/identity_and_access_managing/presentation/http/controllers"
	// identityAndAccessManagingMiddlewares "github.com/z411392/golab/modules/identity_and_access_managing/presentation/http/middlewares"
	systemMaintainingControllers "github.com/z411392/golab/modules/system_maintaining/presentation/http/controllers"
)

func NewHttpHandler() http.Handler {
	router := gin.Default()
	router.GET("/liveness_check", systemMaintainingControllers.OnCheckingLiveness)
	router.GET("/readiness_check", systemMaintainingControllers.OnCheckingReadiness)
	router.POST("/cdc", systemMaintainingControllers.OnChangeDataCaptured)
	router.POST("/webhook", systemMaintainingControllers.Webhook)
	// {
	// 	router := router.Group("/")
	// 	router.Use(identityAndAccessManagingMiddlewares.WithCredentialResolving())
	// 	{
	// 		router := router.Group("/")
	// 		router.Use(identityAndAccessManagingMiddlewares.EnsureUserSignedIn())
	// 		router.GET("/auth/me", identityAndAccessManagingControllers.OnRetrievingProfile)
	// 	}
	// }
	return router
}
