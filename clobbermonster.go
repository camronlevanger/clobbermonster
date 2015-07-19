package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

var (
	protocol              = flag.String("p", "amqp", "use amqp or stomp protocol")
	connection_string     = flag.String("u", "amqp://guest:guest@localhost:5672/", "AMQP Connection String")
	queue_name            = flag.String("n", "queue", "The name of the persistent queue")
	total_messages        = flag.Int("t", 100000, "The total number of messages to send to the queue")
	messages_per_interval = flag.Int("m", 100, "The number of messages to send per interval")
	interval              = flag.Int("i", 1, "The number of seconds to wait before sending messages per interval")
	message_dir           = flag.String("f", "~", "Path of the directory containig the JSON messages to send to the queue")
	stomp_connection      = flag.String("s", "localhost", "STOMP host")
	stomp_user            = flag.String("U", "guest", "STOMP username")
	stomp_pass            = flag.String("P", "guest", "STOMP password")
	stomp_topic           = flag.String("T", "test", "STOMP topic (/topic/your_topic)")
)

var json_messages []map[string]interface{}
var total_sent int = 0

// main entry. parses flags, loads json files, sends to requested publisher
func main() {

	flag.Parse()

	files, err := ioutil.ReadDir(*message_dir)
	failOnError(err, "Couldn't read json message file")

	duration := time.Duration(*interval) * time.Second

	for _, file := range files {
		var message map[string]interface{}
		if strings.HasSuffix(file.Name(), ".json") {
			jfile, err := ioutil.ReadFile(*message_dir + "/" + file.Name())
			failOnError(err, "Can't read file")
			fmt.Printf("Using json file: %s\n", file.Name())
			json.Unmarshal(jfile, &message)
			fmt.Printf("Unmarshalled: %s\n", message)
			json_messages = append(json_messages, message)
		}
	}

	for total_sent < *total_messages {
		fmt.Printf("Preparing to send %d more %s messages\n", *total_messages-total_sent, *protocol)
		var messages []map[string]interface{}
		for i := 0; i < *messages_per_interval; i++ {
			messages = append(messages, json_messages[rand.Intn(len(json_messages))])
		}
		time.Sleep(duration)
		if *protocol == "amqp" {
			PublishAMQP(messages, *connection_string, *queue_name)
		} else if *protocol == "stomp" {
			PublishSTOMP(messages, *stomp_connection, *stomp_user, *stomp_pass, *queue_name, *stomp_topic)
		}
		total_sent = total_sent + len(messages)
	}

}

// function to be called on fatal errors, this kills the app
func failOnError(err error, msg string) {
	if err != nil {
		glog.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
