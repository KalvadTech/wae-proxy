package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"github.com/streadway/amqp"
)

func statping(c *gin.Context) {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")
	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
	secret := c.Param("secret")
	waeProxySecret := os.Getenv("WAE_PROXY_SECRET")
	if secret != waeProxySecret {
		c.String(http.StatusForbidden, "Wrong Secret Key")
		return
	}

	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the connection itself.
	channel, err := connection.Channel()
	
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not Connect to RabbitMQ")
		return
	}
	message := amqp.Publishing{
		Body: []byte("statping_error"),
	}

	// We publish the message to the exahange we created earlier
	err = channel.Publish("wae", "wae", false, false, message)

	if err != nil {
		c.String(http.StatusInternalServerError, "Could not Connect to the Queue")
		return
	}
}
