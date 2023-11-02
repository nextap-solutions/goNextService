package components

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nextap-solutions/goNextService/config"
	"github.com/nextap-solutions/goNextService/metricsz"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type metriczComponent struct {
	config   config.MetricszConfig
	registry *prometheus.Registry

	// Server that's currently being used
	// This used for closing the component
	runningServer *http.Server

	// Channel to close when healhz is disabled
	exitChan chan (bool)
}

func NewMetriczComponent(registry *prometheus.Registry, config config.MetricszConfig) *metriczComponent {
	return &metriczComponent{
		config:   config,
		registry: registry,

		exitChan: make(chan bool, 1),
	}
}

func (hc *metriczComponent) Close(ctx context.Context) error {
	hc.exitChan <- true
	if hc.runningServer == nil {
		return nil
	}

	return hc.runningServer.Shutdown(ctx)
}

func (hc *metriczComponent) Startup() error {
	return nil
}

func (hc *metriczComponent) Run() error {
	if !hc.config.Enabled {
		select {
		case _ = <-hc.exitChan:
			return nil
		}
	}

	api := metricsz.Serve(hc.config, hc.registry)
	hc.runningServer = api

	zap.L().Info(fmt.Sprintf("Starting healhz %s", api.Addr))
	return api.ListenAndServe()
}
