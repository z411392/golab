package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
	Source struct {
		Connector string `json:"connector"`
		Db        string `json:"db"`
		Table     string `json:"table"`
	} `json:"source"`
	Op   string `json:"op"`
	TsMs int    `json:"ts_ms"`
}

func Webhook(ctx *gin.Context) {
	// log.Printf("%v\n", ctx.Request.Header)
	signatureGot := ctx.Request.Header.Get("R-Webhook-Signature")
	if signatureGot == "" {
		ctx.JSON(403, nil)
		return
	}
	requestBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(403, nil)
		return
	}
	defer ctx.Request.Body.Close()
	secret, err := base64.StdEncoding.DecodeString(os.Getenv("WEBHOOK_SIGNING_SECRET"))
	if err != nil {
		ctx.JSON(403, nil)
		return
	}
	h := hmac.New(sha256.New, secret)
	toBeSigned := strings.Join([]string{
		ctx.Request.Header.Get("R-Webhook-Id"),
		ctx.Request.Header.Get("R-Webhook-Timestamp"),
		string(requestBody),
	}, ".")
	// log.Printf("%s\n", toBeSigned)
	h.Write([]byte(toBeSigned))
	signatureExpected := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(signatureGot), []byte(signatureExpected)) {
		ctx.JSON(403, nil)
		return
	}
	ctx.JSON(200, nil)
}
