package healthz

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nextap-solutions/goNextService/config"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Checkable Makes sure the object has the ) function
type Checkable interface {
	Healthz(ctx context.Context) error
}

type CheckableFunc func(ctx context.Context) error

func (f CheckableFunc) Healthz(ctx context.Context) error {
	return f(ctx)
}

// Provider is a provder we can check for healthz
type Provider struct {
	Handle Checkable
	Name   string
}

func NewHealthChecker(providers []Provider) *HealthzChecker {
	return &HealthzChecker{
		Providers: providers,
	}
}

// HealthzChecker contains the instance
type HealthzChecker struct {
	Providers []Provider
}

// returns a http.HandlerFunc for the healthz service
func (h *HealthzChecker) Healthz() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Debug("Handling the Healthz")
		w.Header().Set("Content-Type", "application/json")
		if h.Providers == nil {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("OK"))
			if err != nil {
				zap.L().Error("Internal error:", zap.Error(err))
			}
			return
		}

		response := HealthzResponse{
			Services: []Service{},
			Healthy:  true,
		}

		for _, provider := range h.Providers {
			service := Service{
				Name:    provider.Name,
				Healthy: true,
			}
			err := provider.Handle.Healthz(r.Context())
			if err != nil {
				service.ErrorMessage = err.Error()
				service.Healthy = false
				response.Healthy = false
			}
			response.Services = append(response.Services, service)
		}

		if !response.Healthy {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		jsonBody, err := json.Marshal(response)
		if err != nil {
			zap.L().Error("Unable to marshal errors:", zap.Error(err))
		}

		_, err = w.Write(jsonBody)
		if err != nil {
			zap.L().Error("Internal error:", zap.Error(err))
		}
	})
}

func (h *HealthzChecker) Serve(config config.HealthzConfig) *http.Server {
	router := mux.NewRouter()
	router.Handle("/healthz", h.Healthz())
	router.Handle("/liveliness", Liveness())

	api := http.Server{
		Addr:         "0.0.0.0:" + config.Port,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		Handler:      router,
	}

	return &api
}

func Liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			zap.L().Error("Internal error:", zap.Error(err))
		}
	}
}
