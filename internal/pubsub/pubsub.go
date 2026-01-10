package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType bool

const (
	Durable   SimpleQueueType = true
	Transient SimpleQueueType = false
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

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {

	chnl, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	durable := queueType == Durable
	exclusive := queueType == Transient
	autoDelete := queueType == Transient

	queue, err := chnl.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		chnl.Close()
		return nil, amqp.Queue{}, err
	}

	err = chnl.QueueBind(
		queue.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		chnl.Close()
		return nil, amqp.Queue{}, err
	}

	return chnl, queue, nil
}
