package main

import (
	"fmt"

	gamelogic "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	pubsub "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

var connectionURL = "amqp://guest:guest@127.0.0.1:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Failed to create channel:", err)
		return
	}
	defer ch.Close()

	gamelogic.PrintServerHelp()

	pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, "game_logs", "game_logs.*", pubsub.Durable)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {

		case "pause":
			fmt.Println("Sending pause message")
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				fmt.Println("Failed to publish pause:", err)
			}

		case "resume":
			fmt.Println("Sending resume message")
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				fmt.Println("Failed to publish resume:", err)
			}

		case "quit":
			fmt.Println("Exiting server")
			return

		default:
			fmt.Println("Unknown command:", words[0])
		}
	}
}
