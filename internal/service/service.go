package service

type Service struct {
	db     DBRepo
	rbmq   RabbitMQ
	sender Sender
}

func New(d DBRepo, rq RabbitMQ, s Sender) *Service {
	return &Service{
		db:     d,
		rbmq:   rq,
		sender: s,
	}
}
