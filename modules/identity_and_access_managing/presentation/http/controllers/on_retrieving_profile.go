package controllers

import (
	"github.com/gin-gonic/gin"
	dummyJson "github.com/z411392/golab/adapters/http/dummy_json"
	contextKeys "github.com/z411392/golab/constants/context_keys"
)

func OnRetrievingProfile(ctx *gin.Context) {
	user := ctx.MustGet(contextKeys.User).(*dummyJson.User)
	payload := map[string]interface{}{
		"user": user,
	}
	ctx.JSON(200, map[string]interface{}{"payload": payload})
}
