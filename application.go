package goNextService

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Logger interface {
	Debugf(mgs string, params ...any)
	Infof(mgs string, params ...any)
}

type fmtLogger struct{}

func (l *fmtLogger) Debugf(mgs string, params ...any) {
	fmt.Printf(mgs, params...)
}

func (l *fmtLogger) Infof(mgs string, params ...any) {
	fmt.Printf(mgs, params...)
}

type Component interface {
	// Synchronous startup, container runs this one by one
	Startup() error
	// Asynchronous run,  container runs this in parrallel
	// Has to block until it's done
	// When one of componenst exits this func, container will start shutting down
	Run() error
	// Cleanup function
	Close(ctx context.Context) error
}

var defaultLogger = &fmtLogger{}

type Application struct {
	components []Component
	timeout    time.Duration
	logger     Logger
}

func NewApplications(components ...Component) *Application {
	return &Application{
		components: components,
		timeout:    10 * time.Second,
		logger:     defaultLogger,
	}
}

func (app Application) WithTimeout(timeout time.Duration) Application {
	app.timeout = timeout
	return app
}

func (app Application) WithLogger(logger Logger) Application {
	app.logger = logger

	return app
}

func (app *Application) AddComponent(c Component) {
	if app.components == nil {
		app.components = []Component{}
	}

	app.components = append(app.components, c)
}

func (app *Application) Run() error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	serverErrors := make(chan error, 1)

	for _, component := range app.components {
		component := component
		app.logger.Debugf("Starting component %T", component)
		err := component.Startup()
		if err != nil {
			return fmt.Errorf("Startup error: %w", err)
		}
	}

	for _, component := range app.components {
		component := component
		go func() {
			serverErrors <- component.Run()
		}()
	}

	select {
	case err := <-serverErrors:
		app.logger.Infof("Application Error : %v", err)
		// TODO add shutdown time config
		ctx, cancel := context.WithTimeout(context.Background(), app.timeout)
		defer cancel()
		for _, component := range app.components {
			app.logger.Infof("Closing component %T", component)
			err := component.Close(ctx)
			if err != nil {
				app.logger.Infof("Component %T Closed with error %s", component, err.Error())
				return err
			}
		}
		return err

	case sig := <-shutdown:
		app.logger.Infof("%v : Shuting down gracefully", sig)
		exitChan := make(chan error, 1)
		// TODO add shutdown time config
		ctx, cancel := context.WithTimeout(context.Background(), app.timeout)
		defer cancel()

		go func() {
			var err error
			for _, component := range app.components {
				app.logger.Infof("Closing component %T", component)
				err := component.Close(ctx)
				if err != nil {
					app.logger.Infof("Shuting down did not complete, %v", err)
				}
			}
			exitChan <- err
		}()

		select {
		case err := <-exitChan:
			return err
		case <-ctx.Done():
			return fmt.Errorf("Killed after %s", app.timeout)
		}
	}
}
