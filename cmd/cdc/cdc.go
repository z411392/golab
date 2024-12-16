package cdc

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

const use = "cdc"

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

func handle(message amqp.Delivery) (err error) {
	if string(message.Body) == "default" {
		return
	}
	changeCaptured := &ChangeCaptured{}
	err = json.Unmarshal(message.Body, &changeCaptured)
	if err != nil {
		return
	}
	bytes, err := changeCaptured.Marshal()
	if err != nil {
		return
	}
	log.Printf("%v\n", string(bytes))
	return
}

func runE(command *cobra.Command, args []string) (err error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	connection, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		return
	}
	defer channel.Close()
	messages, err := channel.Consume("cdc", "", false, false, false, false, nil)
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	go (func() {
	loop:
		for {
			select {
			case message := <-messages:
				err := handle(message)
				if err == nil {
					err = message.Ack(false)
				}
				if err == nil {
					continue loop
				}
				message.Nack(false, true)
			case <-ctx.Done():
				break loop
			}
		}
	})()
	<-stop
	cancel()
	return
}

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:  use,
		RunE: runE,
	}
	return command
}
