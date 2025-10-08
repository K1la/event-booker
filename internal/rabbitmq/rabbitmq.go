package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/K1la/event-booker/internal/config"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"time"
)

const (
	retries = 3
)

type RabbitMq struct {
	publisher *rabbitmq.Publisher
	consumer  <-chan amqp.Delivery
}

func New(cfg *config.Config) *RabbitMq {
	publisher, deliveries := initProducerAndConsumer(cfg)

	return &RabbitMq{
		publisher: publisher,
		consumer:  deliveries,
	}
}

func (r *RabbitMq) Publish(booking dto.QueueMessage) error {
	body, err := json.Marshal(booking)
	if err != nil {
		return fmt.Errorf("could not marshal booking to send to rabbitmq:" + err.Error())
	}

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  2,
	}

	headers := amqp.Table{
		"x-delay": (time.Minute * 15).Milliseconds(), // отправляет после 15 минут (период ожидания оплаты)
	}

	options := rabbitmq.PublishingOptions{
		Headers: headers,
	}

	return r.publisher.PublishWithRetry(body, "bookings", "application/json", strategy, options)
}

func (r *RabbitMq) Consume(ctx context.Context) (<-chan []byte, error) {
	messages := make(chan []byte)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(messages)
				return
			default:
				next, ok := <-r.consumer
				if !ok {
					return
				}

				if err := next.Ack(false); err != nil {
					zlog.Logger.Error().Msg("could not acknowledge message consuming: " + err.Error())
				}
				messages <- next.Body
			}
		}
	}()

	return messages, nil
}

func initProducerAndConsumer(cfg *config.Config) (*rabbitmq.Publisher, <-chan amqp.Delivery) {
	url := fmt.Sprintf(
		"amqp://guest:guest@%s%s/",
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	connection, err := rabbitmq.Connect(url, retries, 0)
	if err != nil {
		log.Fatal("could not connect to rabbitmq server: ", err)
	}

	pubCh, err := connection.Channel()
	if err != nil {
		log.Fatal("could not create channel for rabbitmq: ", err)
	}

	args := amqp.Table{"x-delayed-type": "direct"}
	err = pubCh.ExchangeDeclare(
		"bookings",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		log.Fatal("could not declare exchange for rabbitmq: ", err)
	}

	qm := rabbitmq.NewQueueManager(pubCh)
	_, err = qm.DeclareQueue("bookings")
	if err != nil {
		log.Fatal("could not create queue for rabbitmq: ", err)
	}

	err = pubCh.QueueBind("bookings", "bookings", "bookings", false, nil)
	if err != nil {
		log.Fatal("could not bind queue to exchange: ", err)
	}

	publisher := rabbitmq.NewPublisher(pubCh, "bookings")

	conCh, err := connection.Channel()
	if err != nil {
		log.Fatal("could not create channel for consumer rabbitmq: ", err)
	}

	deliveries, err := conCh.Consume("bookings", "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("could not create consumer for rabbitmq: ", err)
	}

	return publisher, deliveries
}
