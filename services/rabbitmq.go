package services

import (
	"log"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

func ConnectRabbitMQ() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	return conn
}

func CreateQueue(conn *amqp.Connection, queueName string) (amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return amqp.Queue{}, err
	}
	defer ch.Close()
	queue, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return amqp.Queue{}, err
	}
	log.Printf("Queue %s created", queue.Name)
	return queue, nil
}

func PublishMessage(conn *amqp.Connection, queueName string, message string) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return err
	}
	defer ch.Close()
	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
		return err
	}
	log.Printf("Message sent to queue %s: %s", queueName, message)
	return nil
}

func ConsumeMessages(conn *amqp.Connection, mongoClient *mongo.Client, database, collection, queueName string) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for m := range msgs {
			log.Printf("Received a message: %s", m.Body)
			err := InsertTodo(mongoClient, database, collection, string(m.Body))
			if err != nil {
				log.Printf("Failed to insert into MongoDB: %v", err)
			} else {
				log.Printf("Message inserted into MongoDB: %s", m.Body)
			}
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}
