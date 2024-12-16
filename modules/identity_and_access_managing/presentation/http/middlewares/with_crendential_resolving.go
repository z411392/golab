package middlewares

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/z411392/golab/adapters/http/dummy_json"
	contextKeys "github.com/z411392/golab/constants/context_keys"
)

func WithCredentialResolving() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		pattern := regexp.MustCompile(`\s+`)
		parts := pattern.Split(authorization, 2)
		if len(parts) != 2 {
			ctx.Next()
			return
		}
		token := parts[1]
		adapter := dummy_json.NewDummyJsonAdapter()
		user, err := adapter.GetAuthUser(token)
		if err != nil {
			ctx.Next()
			return
		}
		if user == nil {
			ctx.Next()
			return
		}
		ctx.Set(contextKeys.User, user)
	}
}
