package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	contextKeys "github.com/z411392/golab/constants/context_keys"
	"github.com/z411392/golab/modules/identity_and_access_managing/errors"
)

func EnsureUserSignedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get(contextKeys.User)
		if !exists || user == nil {
			err := &errors.Unauthorized{}
			ctx.JSON(http.StatusForbidden, err)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
