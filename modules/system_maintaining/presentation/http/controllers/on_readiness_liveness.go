package controllers

import (
	"github.com/gin-gonic/gin"
)

func OnCheckingReadiness(ctx *gin.Context) {
	payload := map[string]interface{}{}
	ctx.JSON(200, map[string]interface{}{
		"payload": payload,
	})
}
