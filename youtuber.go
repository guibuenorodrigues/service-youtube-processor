package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"google.golang.org/api/youtube/v3"
)

var (
	// AMQURL = os.Getenv("AMQ_URL")
	AMQURL = "amqp://guest:guest@127.0.0.1:5672"
)

func sender() {
	// conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	// failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	// ch, err := conn.Channel()
	// failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	// q, err := ch.QueueDeclare(
	// 	"to.processor.postApiData", // name
	// 	true,                       // durable
	// 	false,                      // delete when unused
	// 	false,                      // exclusive
	// 	false,                      // no-wait
	// 	nil,                        // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	// body := "GUILHERME"
	// err = ch.Publish(
	// 	"",     // exchange
	// 	q.Name, // routing key
	// 	false,  // mandatory
	// 	false,
	// 	amqp.Publishing{
	// 		DeliveryMode: amqp.Persistent,
	// 		AppId:        "service.collector",
	// 		ContentType:  "text/plain",
	// 		Body:         []byte(body),
	// 	})
	// failOnError(err, "Failed to publish a message")
	// log.Printf(" [x] Sent %s", body)

	// conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	// failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	// ch, err := conn.Channel()
	// failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	// err = ch.ExchangeDeclare(
	// 	"youtuber", // name
	// 	"direct",   // type
	// 	true,       // durable
	// 	false,      // auto-deleted
	// 	false,      // internal
	// 	false,      // no-wait
	// 	nil,        // arguments
	// )
	// failOnError(err, "Failed to declare an exchange")

	// body := "POOOORRA"
	// err = ch.Publish(
	// 	"youtuber",          // exchange
	// 	"to.processor.post", // routing key
	// 	false,               // mandatory
	// 	false,               // immediate
	// 	amqp.Publishing{
	// 		DeliveryMode: amqp.Persistent,
	// 		AppId:        "service.collector",
	// 		ContentType:  "application/json",
	// 		Body:         []byte(body),
	// 	})
	// failOnError(err, "Failed to publish a message")

	// log.Printf(" [x] Sent %s", body)

}

// MessageResponse from youtuber
type MessageResponse struct {
	LocationName string                    `json:"location"`
	Videos       youtube.VideoListResponse `json:"videos"`
	Interal      interalContentResponse    `json:"internal"`
}

type interalContentResponse struct {
	CorrelationID string `json:"correlationId"`
	AppID         string `json:"appId"`
}

// GetVideoDataFromBroker - listen to retrieve the data from youtebe api broker
func GetVideoDataFromBroker() {

	conn, err := amqp.Dial(AMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		// "to.processor.postApiData", // name
		"to.processor.post",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(1, 0, false)

	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err := youtubeDataReceived(d)

			if err != nil {
				fmt.Println(err)
				fmt.Println("Message not processed, being requeued...")
				d.Reject(false) // not requeue the message
			} else {
				d.Ack(false) // ack the message
			}

		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func youtubeDataReceived(d amqp.Delivery) error {

	// get message and convert to struct

	message := MessageResponse{}

	if string(d.Body) == "" {
		return errors.New("The string received is empty")
	}

	err := json.Unmarshal(d.Body, &message)

	if err != nil {
		return err
	}

	message.Interal.CorrelationID = d.CorrelationId
	message.Interal.AppID = d.AppId

	err = message.PostToProcess()

	if err != nil {
		return err
	}

	return nil
}
