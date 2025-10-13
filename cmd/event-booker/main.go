package main

import (
	"context"
	"github.com/K1la/event-booker/internal/api/handler"
	"github.com/K1la/event-booker/internal/api/router"
	"github.com/K1la/event-booker/internal/api/server"
	"github.com/K1la/event-booker/internal/config"
	"github.com/K1la/event-booker/internal/rabbitmq"
	"github.com/K1la/event-booker/internal/repository"
	"github.com/K1la/event-booker/internal/sender"
	"github.com/K1la/event-booker/internal/service"
	"github.com/wb-go/wbf/zlog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zlog.Init()

	cfg := config.Init()

	db := repository.NewDB(cfg)
	repo := repository.New(db)
	zlog.Logger.Info().Interface("cfg rabbitmq", cfg.RabbitMQ).Msg("cfg rabbitmq in main")
	rabmq := rabbitmq.New(cfg)
	snder := sender.New()
	srvc := service.New(repo, rabmq, snder)

	hndlr := handler.New(srvc)
	r := router.New(hndlr)
	s := server.New(cfg.HTTPServer.Address, r)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// sig channel to handle SIGINT and SIGTERM for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// TODO: use queue in service
	srvc.StartWorker(ctx)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting down signal %v. Shutting down...", sig)
		cancel()
	}()

	if err := s.ListenAndServe(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to start server")
	}
	zlog.Logger.Info().Msg("successfully started server on " + cfg.HTTPServer.Address)
}
