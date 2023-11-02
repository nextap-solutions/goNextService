package application

import "context"

type LifeCycleFunc func(ctx context.Context) error

type lifecycleComponent struct {
	startups []LifeCycleFunc
	run      LifeCycleFunc
	cleanups []LifeCycleFunc

	// Channel to close
	exitChan chan (bool)
}

func NewLifecycleComponent(startups []LifeCycleFunc, run LifeCycleFunc, cleanups []LifeCycleFunc) *lifecycleComponent {
	return &lifecycleComponent{
		startups: startups,
		run:      run,
		cleanups: cleanups,

		exitChan: make(chan bool, 1),
	}
}

func (cc *lifecycleComponent) Startup() error {
	ctx := context.TODO()
	for _, startups := range cc.startups {
		err := startups(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cc *lifecycleComponent) Run() error {
	ctx := context.TODO()
	if cc.run != nil {
		err := cc.run(ctx)
		if err != nil {
			return err
		}
	} else {
		select {
		case _ = <-cc.exitChan:
			return nil
		}
	}
	return nil
}

func (cc *lifecycleComponent) Close(ctx context.Context) error {
	for _, cleanup := range cc.cleanups {
		err := cleanup(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
