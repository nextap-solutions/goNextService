package metricsz

import (
	"net/http"

	"github.com/nextap-solutions/goNextService/config"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(config config.MetricszConfig, reg *prometheus.Registry) *http.Server {
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.InstrumentMetricHandler(
		reg, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	))

	api := http.Server{
		Addr:         "0.0.0.0:" + config.Port,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		Handler:      router,
	}

	return &api
}
