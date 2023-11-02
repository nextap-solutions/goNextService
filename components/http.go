package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/nextap-solutions/goNextService/config"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

type HttpComponent struct {
	httpServer   *http.Server
	serverConfig *config.ServerConfig
	handler      http.Handler
	// Server that's currently being used
	// This used for closing the component
	runningServer *http.Server

	// Channel to close when healhz is disabled
	exitChan chan (bool)
}

type httpComponentOption func(c HttpComponent) HttpComponent

func WithHttpServer(httpServer *http.Server) httpComponentOption {
	return func(c HttpComponent) HttpComponent {
		c.httpServer = httpServer
		return c
	}
}

func WithServerConfig(serverConfig *config.ServerConfig) httpComponentOption {
	return func(c HttpComponent) HttpComponent {
		c.serverConfig = serverConfig
		return c
	}
}

func NewHttpComponent(handler http.Handler, options ...httpComponentOption) *HttpComponent {
	component := HttpComponent{
		handler: handler,

		exitChan: make(chan bool),
	}
	for _, option := range options {
		component = option(component)
	}

	return &component
}

// This is just for testing
func (hc *HttpComponent) GetHandler(t testing.TB) http.Handler {
	t.Helper()
	return hc.handler
}

func (hc *HttpComponent) Startup() error {
	return nil
}

func (hc *HttpComponent) Run() error {
	handler := hc.handler
	if !hc.serverConfig.Enabled {
		select {
		case _ = <-hc.exitChan:
			return nil
		}
	}

	if hc.serverConfig != nil && hc.serverConfig.Cors != nil {
		corsMiddleware := cors.New(cors.Options{
			AllowedOrigins:   hc.serverConfig.Cors.AllowedOrigins,
			AllowCredentials: hc.serverConfig.Cors.AllowCredentials,
			AllowedHeaders:   hc.serverConfig.Cors.AllowedHeaders,
		})
		handler = corsMiddleware.Handler(hc.handler)
	}

	api, err := hc.resolveServer()
	if err != nil {
		return err
	}
	hc.runningServer = api
	api.Handler = handler

	zap.L().Info(fmt.Sprintf("Starting server %s", api.Addr))
	return api.ListenAndServe()
}

func (hc *HttpComponent) Close(ctx context.Context) error {
	if !hc.serverConfig.Enabled {
		hc.exitChan <- true
		close(hc.exitChan)
		return nil
	}
	if hc.runningServer == nil {
		return nil
	}
	err := hc.runningServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	hc.runningServer = nil
	return nil
}

func (hc *HttpComponent) resolveServer() (*http.Server, error) {
	if hc.httpServer != nil && hc.serverConfig != nil {
		return nil, errors.New("Cannot specify both httpServer and serverConfig")
	}
	if hc.serverConfig != nil {
		return &http.Server{
			Addr:              "0.0.0.0:" + hc.serverConfig.Port,
			ReadHeaderTimeout: hc.serverConfig.ReadTimeout,
			WriteTimeout:      hc.serverConfig.WriteTimeout,
		}, nil
	}
	if hc.httpServer != nil {
		return hc.httpServer, nil
	}

	return nil, errors.New("Cannot one of httpServer or serverConfig")
}
