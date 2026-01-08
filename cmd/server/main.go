package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pubsub "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

var connection string = "amqp://guest:guest@127.0.0.1:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	con, err := amqp.Dial(connection)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer con.Close()

	chnl, err := con.Channel()
	if err != nil {
		fmt.Println("Failed to create channel", err)
		return
	}

	fmt.Println("connection successful")

	pubsub.PublishJSON(chnl, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("closing program")
}
