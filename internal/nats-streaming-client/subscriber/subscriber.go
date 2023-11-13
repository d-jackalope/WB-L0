package subscriber

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/d-jackalope/L0/internal/app"
	"github.com/d-jackalope/L0/pkg/models"
	"github.com/nats-io/stan.go"
)

var ErrOrderIsEmpty = errors.New("order_uid is empty")

var sub *Subscriber

type Subscriber struct {
	App      *app.Application
	Subject  string
	ErrChan  chan error
	DoneChan chan struct{}
}

func NewSubscriber(app *app.Application, subject string) *Subscriber {
	if sub != nil {
		return sub
	}
	sub = &Subscriber{
		App:      app,
		Subject:  subject,
		ErrChan:  make(chan error),
		DoneChan: make(chan struct{}),
	}
	return sub
}

func GetSubscriber() Subscriber {
	if sub != nil {
		return *sub
	}
	return Subscriber{}
}

func (s *Subscriber) Run(ctx context.Context) {
	app := app.GetApp()
	_, err := s.App.Sc.Subscribe(s.Subject, msgHandler, stan.DurableName(s.Subject))
	if err != nil {
		s.ErrChan <- err
	}
	app.Log.Info("The app is subscribed to the \"%s\" channel", s.Subject)

	<-ctx.Done() //ожидание сигнала завершения программы

	if err = s.App.Sc.Close(); err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.Log.Output(2, trace)
	} else {
		app.Log.Info("The Nats-streaming connection was successfully closed")
	}
	close(s.DoneChan)
}

func msgHandler(msg *stan.Msg) {
	app := app.GetApp()

	data := msg.Data
	order := models.Order{}
	err := json.Unmarshal(data, &order)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.Log.Output(2, trace)
		return
	}

	if order.OrderUID == "" {
		trace := fmt.Sprintf("%s\n%s", ErrOrderIsEmpty.Error(), debug.Stack())
		app.Log.Output(2, trace)
		return
	}

	app.Cache.Set(order.OrderUID, order)

	exist, err := app.Postgres.Orders.Exist(order.OrderUID)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.Log.Output(2, trace)
	}

	if !exist {
		if err := app.Postgres.Orders.Insert(order.OrderUID, data); err != nil {
			trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
			app.Log.Output(2, trace)
		}
	}

}
