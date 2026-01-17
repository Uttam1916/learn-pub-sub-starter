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
		pubsub.Transient)
	if err != nil {
		fmt.Println("Failed to declare and bind queue:", err)
		return
	}
	defer ch.Close()

	gmeste := gamelogic.NewGameState(username)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {

		case "spawn":
			fmt.Println("Spawning")
			err := gmeste.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}

		case "move":
			fmt.Println("Sending resume message")
			_, err := gmeste.CommandMove(words)
			if err != nil {
				fmt.Println(err)
			}

		case "help":
			gamelogic.PrintClientHelp()

		case "status":
			gmeste.CommandStatus()

		case "spam":
			fmt.Println("Spamming not allowed yet!")

		case "quit":
			fmt.Println("Exiting client")
			return

		default:
			fmt.Println("Unknown command:", words[0])
		}
	}

}
