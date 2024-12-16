package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/z411392/golab/container"
)

type ChangeCaptured struct {
	Payload struct {
		Before interface{} `json:"before"`
		After  interface{} `json:"after"`
		Source struct {
			Connector string `json:"connector"`
			Db        string `json:"db"`
			Table     string `json:"table"`
		} `json:"source"`
		Op   string `json:"op"`
		TsMs int    `json:"ts_ms"`
	} `json:"payload"`
}

func (changeCaptured ChangeCaptured) Marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"before": changeCaptured.Payload.Before,
		"after":  changeCaptured.Payload.After,
		"source": changeCaptured.Payload.Source,
		"op":     changeCaptured.Payload.Op,
		"ts_ms":  changeCaptured.Payload.TsMs,
	})
}

var (
	queueCreated *string
)

func OnChangeDataCaptured(ctx *gin.Context) {
	// log.Printf("%v\n", ctx.Request.Header)
	segments := strings.Split(ctx.Request.Header.Get("Webhook-Signature"), ",")
	if len(segments) < 2 {
		ctx.JSON(403, nil)
		return
	}
	signatureGot := segments[1]
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
	h.Write([]byte(strings.Join([]string{
		ctx.Request.Header.Get("Webhook-Id"),
		ctx.Request.Header.Get("Webhook-Timestamp"),
		string(requestBody),
	}, ".")))
	signatureExpected := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(signatureGot), []byte(signatureExpected)) {
		ctx.JSON(403, nil)
		return
	}
	changeDataCaptured := ChangeCaptured{}
	err = json.Unmarshal(requestBody, &changeDataCaptured)
	if err != nil {
		ctx.JSON(403, nil)
	}
	responseBody, err := changeDataCaptured.Marshal()
	if err != nil {
		ctx.JSON(403, nil)
	}
	h = hmac.New(sha256.New, secret)
	toBeSigned := strings.Join([]string{
		ctx.Request.Header.Get("Webhook-Id"),
		ctx.Request.Header.Get("Webhook-Timestamp"),
		string(responseBody),
	}, ".")
	// log.Printf("%s\n", toBeSigned)
	h.Write([]byte(toBeSigned))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	err = container.Container.Invoke(func(channel *amqp.Channel) error {
		if queueCreated == nil {
			queue, err := channel.QueueDeclare("webhook", true, false, false, false, nil)
			if err != nil {
				return err
			}
			queueCreated = &queue.Name
		}
		err := channel.Publish("", *queueCreated, true, false, amqp.Publishing{
			Headers: amqp.Table{
				"Webhook-Id":        ctx.Request.Header.Get("Webhook-Id"),
				"Webhook-Timestamp": ctx.Request.Header.Get("Webhook-Timestamp"),
				"Webhook-Signature": signature,
			},
			Body: responseBody,
		})
		return err
	})
	if err != nil {
		ctx.JSON(403, nil)
	}
	ctx.JSON(200, nil)
}
