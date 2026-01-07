package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	fmt.Println("connection successful")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("closing program")
}
