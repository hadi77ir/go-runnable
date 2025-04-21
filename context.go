package runnable

import (
	"context"
	"sync"
	"time"

	"github.com/hadi77ir/go-logging"
)

const ContextValuesOffset = 0xFF0
const (
	ContextValueKeyWaitGroup = ContextValuesOffset + iota
	ContextValueKeyLogger
	ContextValueKeyConfig
)

func ContextConfig(c context.Context) any {
	return c.Value(ContextValueKeyConfig)
}

func ContextLogger(c context.Context) logging.Logger {
	l, ok := c.Value(ContextValueKeyLogger).(logging.Logger)
	if !ok || l == nil {
		return logging.NoOpLogger(0)
	}
	return l
}

func ContextWaitGroup(c context.Context) *sync.WaitGroup {
	return c.Value(ContextValueKeyWaitGroup).(*sync.WaitGroup)
}

type valuesCtx struct {
	parent context.Context
	offset int
	val    []any
}

func (c *valuesCtx) Deadline() (deadline time.Time, ok bool) {
	deadline, ok = c.parent.Deadline()
	return
}

func (c *valuesCtx) Done() <-chan struct{} {
	return c.parent.Done()
}

func (c *valuesCtx) Err() error {
	return c.parent.Err()
}

func (c *valuesCtx) Value(key any) any {
	if keyInt, ok := key.(int); ok {
		keyInt -= c.offset
		if keyInt >= 0 && keyInt < len(c.val) {
			return c.val[keyInt]
		}
	}

	return c.parent.Value(key)
}

// ContextWithValues creates a context with given values stored as an array of
// sequentially indexed values. Values may be accessed with the seq index.
// An offset can be used to prevent accidental overrides.
func ContextWithValues(parent context.Context, offset int, values ...any) context.Context {
	return &valuesCtx{
		parent: parent,
		offset: offset,
		val:    values,
	}
}
