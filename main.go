package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"github.com/streadway/amqp"
	"io/ioutil"
)

func main() {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")
	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)
	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
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
	router.POST("/webhook/:secret", func(c *gin.Context) {
		secret := c.Param("secret")
		waeProxySecret := os.Getenv("WAE_PROXY_SECRET")
		if secret != waeProxySecret {
			c.String(http.StatusForbidden, "Wrong Secret Key")
			return
		}
		var bodyBytes []byte
		var err error
		if c.Request.Body != nil {
			bodyBytes, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.String(http.StatusInternalServerError, "No Body")
				return
			}
		}

		// Create a channel from the connection. We'll use channels to access the data in the queue rather than the connection itself.
		channel, err := connection.Channel()
		
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not Connect to RabbitMQ")
			return
		}
		message := amqp.Publishing{
			Body: bodyBytes,
		}

		// We publish the message to the exahange we created earlier
		err = channel.Publish("wae", "wae", false, false, message)

		if err != nil {
			c.String(http.StatusInternalServerError, "Could not Connect to the Queue")
			return
		}
	})
	router.Run(":8080")
}
