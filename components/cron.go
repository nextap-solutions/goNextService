package components

import (
	"context"

	"github.com/go-co-op/gocron/v2"
)

type cronComponent struct {
	cron gocron.Scheduler
}

func NewCronComponent(cron gocron.Scheduler) *cronComponent {
	return &cronComponent{
		cron: cron,
	}
}

func (cc *cronComponent) Startup() error {
	return nil
}

func (cc *cronComponent) Run() error {
	cc.cron.Start()
	select {}
}

func (cc *cronComponent) Close(ctx context.Context) error {
	return cc.cron.Shutdown()
}
