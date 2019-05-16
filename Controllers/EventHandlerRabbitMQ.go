package Controllers

import (
	"TaiBaiSupport/Models"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var RabbitmqEventReceivedChan chan Models.TaibaiClassroomEvent
var RabbitmqEventTobeSendChan chan Models.TaibaiClassroomEvent
var ExchangeName = "taibai-exchange"

func init()  {
	RabbitmqEventReceivedChan = make(chan Models.TaibaiClassroomEvent)
	RabbitmqEventTobeSendChan = make(chan Models.TaibaiClassroomEvent)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to get a channel")

	err = ch.ExchangeDeclare(
		ExchangeName,   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare exchange")

	hostName,_ := os.Hostname()
	q, err := ch.QueueDeclare(
		"taibai-queue-"+hostName,    // name
		false, // durable
		true, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare queue")


	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		ExchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind exchange with queue")

	consumer, err := ch.Consume(
		q.Name, // queue
		"taibai-consumer-"+hostName,     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")


	go func() {
		for event := range consumer {
			taibaiEvent := Models.TaibaiClassroomEvent{}
			err := json.Unmarshal(event.Body, &taibaiEvent)
			if err!=nil{
				log.Printf("failed to Unmarshal to TaibaiEvent message: %s" , event.Body)
			}else {
				RabbitmqEventReceivedChan <- taibaiEvent
			}
		}
	}()

	go func() {
		for event := range RabbitmqEventTobeSendChan {
			message, _ := json.Marshal(event)
			err = ch.Publish(
				ExchangeName, // exchange
				"",     // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        message,
				})
		}
	}()

}
