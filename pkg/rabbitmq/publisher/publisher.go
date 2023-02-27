package publisher

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/urosradivojevic/health/pkg/rabbitmq/messages"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Publish(ID string, queueName string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare(queueName, true, false, false, false, nil)
	handleError(err, "Could not declare queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure Qos")

	campaign := &messages.Campaign{
		ID: ID,
	}

	body, err := json.Marshal(campaign)
	handleError(err, "Error encoding JSON")
	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		log.Fatal("Error publishing message: %s", err)
	}
	log.Printf("ID: %v", campaign.ID)
}
