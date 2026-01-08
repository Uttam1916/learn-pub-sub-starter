package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	json, err := json.Marshal(val)
	if err != nil {
		fmt.Printf("error marshaling value, err: %v", err)
		return err
	}

	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: json})
	return nil
}
