package go_runnable

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hadi77ir/go-logging"
)

var ErrStillRunning = errors.New("runnable running")
var ErrAlreadyRunning = ErrStillRunning
var ErrRunFuncNil = errors.New("run func is nil")

type Runnable func(ctx context.Context) error

func Run(r Runnable, config any, logger logging.Logger, ctx context.Context) error {
	wg := &sync.WaitGroup{}
	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithCancel(ctx)
	// for order of arguments, check ContextValueKey* in context.go
	ctx = ContextWithValues(ctx, ContextValuesOffset, wg, logger, config)

	wg.Add(1)
	errChan := make(chan error, 1)
	go func() {
		defer wg.Done()
		errChan <- r(ctx)
		cancelFunc()
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	logger.Log(logging.InfoLevel, "stopping")

	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-ch:
	case <-time.After(time.Duration(10) * time.Second):
		logger.Log(logging.WarnLevel, "some goroutines have not stopped yet")
	}

	select {
	case err := <-errChan:
		return err
	default:
	}
	return nil
}
