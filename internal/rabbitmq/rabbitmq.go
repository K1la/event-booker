package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/K1la/event-booker/internal/config"
	"github.com/K1la/event-booker/internal/dto"
	"log"
	"os"

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
	publisher, deliveries, err := initProducerAndConsumer(cfg)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to initialize RabbitMQ, running without queue functionality")
		return &RabbitMq{
			publisher: nil,
			consumer:  nil,
		}
	}

	return &RabbitMq{
		publisher: publisher,
		consumer:  deliveries,
	}
}

func (r *RabbitMq) Publish(booking dto.QueueMessage) error {
	if r.publisher == nil {
		return fmt.Errorf("rabbitmq publisher is not initialized")
	}

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
	if r.consumer == nil {
		return nil, fmt.Errorf("rabbitmq consumer is not initialized")
	}

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

func initProducerAndConsumer(cfg *config.Config) (*rabbitmq.Publisher, <-chan amqp.Delivery, error) {
	zlog.Logger.Info().Interface("cfg", cfg.RabbitMQ).Msg("cfg rabbitmq in rabbitmq")
	var host, port string
	if cfg.RabbitMQ.Port == "" || cfg.RabbitMQ.Host == "" {
		host = os.Getenv("RABBIT_HOST")
		port = os.Getenv("RABBIT_PORT")
	} else {
		host = cfg.RabbitMQ.Host
		port = cfg.RabbitMQ.Port
	}
	url := fmt.Sprintf(
		"amqp://guest:guest@%s%s/",
		host,
		port,
	)
	log.Printf("connecting to rabbitmq url=[ %s ]", url)

	zlog.Logger.Info().Str("url", url).Msg("connecting to rabbitmq")
	connection, err := rabbitmq.Connect(url, retries, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to rabbitmq server: %w", err)
	}

	pubCh, err := connection.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("could not create channel for rabbitmq: %w", err)
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
		return nil, nil, fmt.Errorf("could not declare exchange for rabbitmq: %w", err)
	}

	qm := rabbitmq.NewQueueManager(pubCh)
	_, err = qm.DeclareQueue("bookings")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create queue for rabbitmq: %w", err)
	}

	err = pubCh.QueueBind("bookings", "bookings", "bookings", false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not bind queue to exchange: %w", err)
	}

	publisher := rabbitmq.NewPublisher(pubCh, "bookings")

	conCh, err := connection.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("could not create channel for consumer rabbitmq: %w", err)
	}

	deliveries, err := conCh.Consume("bookings", "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create consumer for rabbitmq: %w", err)
	}

	return publisher, deliveries, nil
}
