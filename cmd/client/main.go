package main

import (
	"fmt"
	"os"
	"os/signal"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pubsub "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

var connectionURL = "amqp://guest:guest@127.0.0.1:5672/"

func main() {
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		fmt.Println("Could not connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Could not read username:", err)
		return
	}

	queueName := routing.PauseKey + "." + username

	ch, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		fmt.Println("Failed to declare and bind queue:", err)
		return
	}
	defer ch.Close()

	// Wait for Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
