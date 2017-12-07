package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal(err)
		panic(msg)
	}
}

type server struct {
	Url      string
	Username string
	Password string
}

func (s *server) init() {
	//read config
	file, _ := os.Open("server.conf")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(s)
	failOnError(err, "Decode error")
	log.Println("DB URL:", s.Url)
	log.Println("DB Username:", s.Username)
	log.Println("DB Password:", s.Password)
	file.Close()

}

func main() {
	config := server{}
	config.init()
	conn, err := amqp.Dial("amqp://" + config.Username + ":" + config.Password + "@" + config.Url)
	failOnError(err, "Error connect to rabbitmq")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	for d := range msgs {
		log.Println("Received a mesage: ", string(d.Body))
	}

}
