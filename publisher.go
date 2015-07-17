package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/streadway/amqp"
	"gopkg.in/stomp.v1"
)

// Send the messages via AMQP protocol
func PublishAMQP(messages []map[string]interface{}, connection string, queue string) {

	var err error

	glog.Infof("Connecting to queue at %s\n", connection)
	conn, err := amqp.Dial(connection)
	failOnError(err, "Can't connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Can't get channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		failOnError(err, "Can't declare queue")
	}

	for _, message := range messages {

		msg, _ := json.Marshal(message)

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})
		if err != nil {
			failOnError(err, "Unable to publish, shutting down...")
		}
	}

}

// Send the messages via STOMP protocol
func PublishSTOMP(messages []map[string]interface{}, connection string, user string, pass string, queue string, topic string) {

	conn, err := stomp.Dial("tcp", connection, stomp.Options{
		Login:    user,
		Passcode: pass,
		Host:     "/",
	})
	if err != nil {
		glog.Errorf("Unable to connect to STOMP server: %s\n", err)
		return
	}

	for _, message := range messages {

		encoded_message, _ := json.Marshal(message)

		err = conn.Send(
			topic,
			"text/plain",
			encoded_message,
			nil,
		)
		if err != nil {
			glog.Errorf("Error sending STOMP message: %s : %s\n", err, encoded_message)
			conn.Disconnect()
			return
		}
	}

	conn.Disconnect()

}
