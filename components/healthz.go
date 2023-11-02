package components

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nextap-solutions/goNextService/config"
	"github.com/nextap-solutions/goNextService/healthz"

	"go.uber.org/zap"
)

type healthzComponent struct {
	config    config.HealthzConfig
	providers []healthz.Provider

	// Server that's currently being used
	// This used for closing the component
	runningServer *http.Server

	// Channel to close when healhz is disabled
	exitChan chan (bool)
}

func NewHealthzComponent(providers []healthz.Provider, config config.HealthzConfig) *healthzComponent {
	return &healthzComponent{
		config:    config,
		providers: providers,

		exitChan: make(chan bool, 1),
	}
}

func (hc *healthzComponent) Close(ctx context.Context) error {
	hc.exitChan <- true
	if hc.runningServer == nil {
		return nil
	}

	return hc.runningServer.Shutdown(ctx)
}

func (hc *healthzComponent) Startup() error {
	return nil
}

func (hc *healthzComponent) Run() error {
	if !hc.config.Enabled {
		select {
		case _ = <-hc.exitChan:
			return nil
		}
	}

	healthzChecker := healthz.NewHealthChecker(hc.providers)
	api := healthzChecker.Serve(hc.config)
	hc.runningServer = api

	zap.L().Info(fmt.Sprintf("Starting healhz %s", api.Addr))
	return api.ListenAndServe()
}
