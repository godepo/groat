package ctxgroup

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithWaitGroup(t *testing.T) {
	wg := &sync.WaitGroup{}
	ctx := WithWaitGroup(context.Background(), wg)
	ch := make(chan struct{})
	IncAt(ctx)

	go func() {
		<-ch
		DoneFrom(ctx)
	}()

	go func() {
		time.Sleep(time.Millisecond * 10)
		ch <- struct{}{}
	}()

	wg.Wait()
}

func TestDoneFrom(t *testing.T) {
	t.Run("should be able to return false at empty context", func(t *testing.T) {
		assert.False(t, DoneFrom(context.Background()))
	})
}

func TestIncAt(t *testing.T) {
	t.Run("should be able to false at empty context", func(t *testing.T) {
		assert.False(t, IncAt(context.Background()))
	})
}
