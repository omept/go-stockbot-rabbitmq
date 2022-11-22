package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ong-gtp/go-stockbot/api"
	"github.com/ong-gtp/go-stockbot/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type StockRequest struct {
	ChatRoomName string `json:"chatRoomName"`
	ChatRoomId   uint   `json:"chatRoomId"`
	ChatMessage  string `json:"chatMessage"`
}

type StockReponse struct {
	RoomId  uint   `json:"RoomId"`
	Message string `json:"Message"`
}

type Broker struct {
	ReceiverQueue  amqp.Queue
	PublisherQueue amqp.Queue
	Channel        *amqp.Channel
}

func (b *Broker) SetUp(ch *amqp.Channel) {
	receiverQueue := os.Getenv("STKBT_RECEIVER_QUEUE")
	publisherQueue := os.Getenv("STKBT_PUBLISHER_QUEUE")

	q1, err := ch.QueueDeclare(
		receiverQueue, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	utils.FailOnError(err, "Failed to declare"+receiverQueue+" queue")

	q2, err := ch.QueueDeclare(
		publisherQueue, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	utils.FailOnError(err, "Failed to declare "+publisherQueue+" queue")

	b.ReceiverQueue = q1
	b.PublisherQueue = q2
	b.Channel = ch
}

func (b *Broker) PublishMessage(sr StockReponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(sr)
	if err != nil {
		log.Printf("Response structure error %s ", err)
	}

	err = b.Channel.PublishWithContext(ctx,
		"",                    // exchange
		b.PublisherQueue.Name, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	utils.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func (b *Broker) ReadMessages() {
	msgs, err := b.Channel.Consume(
		b.ReceiverQueue.Name, // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	rsvdMsgs := make(chan StockRequest)
	go messageTransformer(msgs, rsvdMsgs)
	go processRequest(rsvdMsgs, b)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}

func messageTransformer(entries <-chan amqp.Delivery, receivedMessages chan StockRequest) {
	var sr StockRequest
	for d := range entries {
		log.Println("d.Body", string(d.Body))
		err := json.Unmarshal([]byte(d.Body), &sr)
		if err != nil {
			log.Printf("Error on received request : %s ", err)
			continue
		}
		log.Println("Received a request")
		receivedMessages <- sr
	}
}

func processRequest(s <-chan StockRequest, b *Broker) {

	for r := range s {
		log.Println("processing stock request for ", r.ChatRoomId)
		cM := r.ChatMessage
		cM = strings.Replace(cM, "/stock=", "", 1)
		sr := StockReponse{
			RoomId:  r.ChatRoomId,
			Message: fmt.Sprintf("Processing: %s", cM),
		}
		go b.PublishMessage(sr)
		msg := api.EvalStock(cM)
		sr2 := StockReponse{
			RoomId:  r.ChatRoomId,
			Message: msg,
		}
		go b.PublishMessage(sr2)
		log.Println("processed", r.ChatMessage)
	}
}
