package main

import (
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"os"
)

func main() {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")
	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
	defer connection.Close()
	channel, err := connection.Channel()
	defer channel.Close()
	durable, exclusive := true, false
	autoDelete, noWait := false, false
	q, err := channel.QueueDeclare("wae-light", durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		panic("could not create queue in RabbitMQ:" + err.Error())
	}
	channel.QueueBind(q.Name, "#", "wae", false, nil)
	router := gin.Default()
	router.POST("/webhook/clevercloud/:secret", clevercloud)
	router.POST("/webhook/statping/:secret", statping)
	router.Run(":8080")
}
