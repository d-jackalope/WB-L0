package main

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/d-jackalope/L0/internal/app"
	"github.com/d-jackalope/L0/internal/config"
	"github.com/d-jackalope/L0/internal/nats-streaming-client/subscriber"
	"github.com/d-jackalope/L0/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/stan.go"
)

func main() {
	if err := config.ReadConfigYAML("config.yaml"); err != nil {
		log.Fatalf("parse config:  %v", err)
	}

	cfg := config.GetConfig()

	ctx, cancel := context.WithCancel(context.Background())

	pool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("error connection to postgre:  %v", err)
	}
	defer pool.Close()

	sc, err := stan.Connect(cfg.NatsStreaming.ClusterID, cfg.NatsStreaming.ClientID, stan.NatsURL(cfg.NatsStreaming.URL))
	if err != nil {
		log.Fatalf("error connecting to NATS Streaming: %v", err)
	}

	app := app.NewApp(ctx, cfg.Server.URL, pool, sc)

	if err = app.CreateNewTablesInPostgre(); err != nil {
		log.Fatalf("create table: %v", err)
	}

	if err := app.GetCacheFromPostgre(); err != nil {
		log.Fatalf("get cache: %v", err)
	}

	subscriber := subscriber.NewSubscriber(&app, cfg.NatsStreaming.Subject)
	// закрытие подключений происходит в горутине после сигнала об отключении приложения
	go subscriber.Run(ctx)
	server := server.New(&app)
	go server.Run(ctx)

	go app.SoftShutdown(cancel)
	go service(&app, subscriber, server, cancel)

	<-app.Ctx.Done()
	app.Log.Info("Waiting for goroutines to close...")

	<-subscriber.DoneChan
	<-server.DoneChan
	app.Log.Info("Done.")
}

func service(app *app.Application, sub *subscriber.Subscriber, server *server.Server, cansel context.CancelFunc) {
	select {
	case err := <-sub.ErrChan:
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.Log.Output(2, trace)
		cansel()
	case err := <-server.ErrChan:
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.Log.Output(2, trace)
		cansel()

	}
}
