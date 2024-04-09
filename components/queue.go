package components

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type QueueHandler func(chan error) error

type queueComponent struct {
	handlers []QueueHandler
	close    func(ctx context.Context) error
}

func NewQueueComponent(handlers []QueueHandler, options ...queueComponentOption) *queueComponent {
	component := queueComponent{
		handlers: handlers,
	}
	for _, option := range options {
		component = option(component)
	}

	return &component
}

type queueComponentOption func(q queueComponent) queueComponent

func WithQueueClose(close func(ctx context.Context) error) queueComponentOption {
	return func(q queueComponent) queueComponent {
		q.close = close
		return q
	}
}

func (qc *queueComponent) Close(ctx context.Context) error {
	if qc.close == nil {
		return nil
	}
	return qc.close(ctx)
}

func (qc *queueComponent) Startup() error {
	return nil
}

func (qc *queueComponent) Run() error {
	zap.L().Info(fmt.Sprintf("Starting queue component with %d", len(qc.handlers)))
	queueHandlerErrors := make(chan error, 10)
	queueErrors := make(chan error, 1)

	for _, handler := range qc.handlers {
		handler := handler
		go func() {
			queueErrors <- handler(queueHandlerErrors)
		}()
	}
	go func() {
		for err := range queueHandlerErrors {
			zap.L().Error("Error from consumer queue", zap.Error(err))
		}
	}()

	err := <-queueErrors
	return err
}
