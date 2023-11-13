package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/d-jackalope/L0/internal/cache"
	"github.com/d-jackalope/L0/internal/db"
	"github.com/d-jackalope/L0/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/stan.go"
)

var app *Application

type Application struct {
	Ctx      context.Context
	Srv      *http.Server
	Sc       stan.Conn
	Log      logger.Logger
	Cache    cache.Cache
	Postgres *db.Postgres
}

func NewApp(ctx context.Context, addr string, pool *pgxpool.Pool, sc stan.Conn) Application {
	if app != nil {
		return *app
	}

	app = &Application{
		Ctx:      ctx,
		Sc:       sc,
		Log:      logger.New(),
		Cache:    cache.New(),
		Postgres: db.NewPostgres(pool, ctx),
	}

	app.Srv = &http.Server{
		Addr:     addr,
		ErrorLog: app.Log.Err(),
	}

	return *app
}

func GetApp() Application {
	if app != nil {
		return *app
	}
	return Application{}
}

func (app *Application) CreateNewTablesInPostgre() error {
	if err := app.Postgres.Orders.Create(); err != nil {
		return err
	}
	return nil
}

func (app *Application) GetCacheFromPostgre() error {
	cache, err := app.Postgres.Orders.GetAllData() //Возвращает пустую мапу в случае, если кеша еще нет в таблице
	if err != nil {
		return err
	}

	if len(cache) > 0 {
		app.Cache.Update(cache) //обновляет кеш данными из Potsgres
		app.Log.Info("Сache received successfully")
	}

	return nil
}

func (app *Application) SoftShutdown(canсel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	app.Log.Info("The application is running. Press Ctrl+C (SIGINT) to end")
	<-c
	app.Log.Info("Closing the application...")
	canсel() //сигнал о завершении работы приложения
}
