package integration_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/streadway/amqp"
	"github.com/z411392/golab/container"
)

func Test_要能與AmqpServer建立連線(t *testing.T) {
	t.SkipNow()
	err := container.Container.Invoke(func(channel *amqp.Channel) (err error) {
		queue, err := channel.QueueDeclare("test", true, false, false, false, nil)
		if err != nil {
			return
		}
		body, _ := json.Marshal(map[string]interface{}{"a": "b"})
		hmac := hmac.New(sha256.New, []byte(os.Getenv("WEBHOOK_SIGNING_SECRET")))
		hmac.Write(body)
		signature := hex.EncodeToString(hmac.Sum(nil))
		err = channel.Publish("", queue.Name, true, false, amqp.Publishing{
			Headers: amqp.Table{
				"Signature": signature,
			},
			Body: body,
		})
		return
	})
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}
}
