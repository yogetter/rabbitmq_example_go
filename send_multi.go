package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
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

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := argsPass(os.Args)
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Print(" [x] Sent ", body)
	failOnError(err, "Failed to publish a message")
}

func argsPass(args []string) string {
	if (len(args) < 2) || os.Args[1] == "" {
		return "Hello"
	} else {
		return strings.Join(args[1:], " ")
	}
}
