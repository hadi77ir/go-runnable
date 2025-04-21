package runnable

import (
	"context"
	"sync"
	"testing"

	"github.com/hadi77ir/go-logging"
	"github.com/hadi77ir/go-logging/logrus"
)

func TestContext(t *testing.T) {
	c := ContextWithValues(context.Background(), 10, 0, 1, 2, 3, 4, 5)
	for i := 0; i < 6; i++ {
		if c.Value(10+i) != i {
			t.FailNow()
		}
	}
}
func TestContextLogger(t *testing.T) {
	l, _ := logrus.New("a")
	c := ContextWithValues(context.Background(), ContextValuesOffset, &sync.WaitGroup{}, l, nil)
	newL := ContextLogger(c)
	if newL == logging.NoOpLogger(0) {
		t.FailNow()
	}
}
