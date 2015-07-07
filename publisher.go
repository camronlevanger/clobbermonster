package main

import (
    "encoding/json"
    "github.com/golang/glog"
    "github.com/streadway/amqp"
)

func Publish(messages []map[string]interface{}, connection string, queue string) {

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
