package Controllers

import (
	"TaibaiSupport/Models"
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


// 不断消费mq的消息


var RabbitmqEventReceivedChan chan Models.TaibaiClassroomEvent
var RabbitmqEventTobeSendChan chan Models.TaibaiClassroomEvent
var ExchangeName = "taibai-exchange-"	// + region
var QueueName = "taibai-queue-"			// + region + hostname
var ConsumerName = "taibai-consumer-"	// + region + hostname

func init()  {

	RabbitmqEventReceivedChan = make(chan Models.TaibaiClassroomEvent, 3)
	RabbitmqEventTobeSendChan = make(chan Models.TaibaiClassroomEvent, 3)

	hostName,_ := os.Hostname()
	rabbitmq_addr := os.Getenv("rabbitmq_addr")
	rabbitmq_user := os.Getenv("rabbitmq_user")
	rabbitmq_passwd := os.Getenv("rabbitmq_passwd")
	classroom_region := os.Getenv("classroom_region")

	ExchangeName = ExchangeName + classroom_region
	QueueName = QueueName + classroom_region + "-" +hostName
	ConsumerName = ConsumerName + classroom_region  + "-" + hostName

	amqp_link := "amqp://" + rabbitmq_user + ":" + rabbitmq_passwd + "@" + rabbitmq_addr +":5672/"
	log.Println(amqp_link)
	conn, err := amqp.Dial(amqp_link)
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

	q, err := ch.QueueDeclare(
		QueueName,    // name
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
		QueueName, // queue
		ConsumerName,     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")


	// 此协程将mq消息转为event存起来
	go func() {
		log.Println("start listen to rabbitmq")
		for event := range consumer {
			taibaiEvent := Models.TaibaiClassroomEvent{}
			err := json.Unmarshal(event.Body, &taibaiEvent)
			log.Println("receive origin mq message:", string(event.Body))
			if err!=nil{
				log.Printf("failed to Unmarshal to TaibaiEvent message: %s" , event.Body)
			}else {
				RabbitmqEventReceivedChan <- taibaiEvent
			}
		}
	}()

	// 此协程消费mq的event
	go func(){
		for event := range RabbitmqEventReceivedChan {
			eventJson,_ := json.Marshal(event)
			log.Println("从mq收到：", string(eventJson))
			switch event.EventType {
			case Models.EventType_UserOnlineStatusChangd:
				HandleEventUserOnlineStatusChanged(&event)
			case Models.EventType_UserVideoPositionChanged:
				HandleEventUserVideoPositionChanged(&event)
			case Models.EventType_1V1StateChanged:
				HandleEvent1V1StateChanged(&event)
			}
		}
	}()

	// 此协程不断给mq发送消息
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
