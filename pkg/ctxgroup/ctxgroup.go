package ctxgroup

import (
	"context"
	"sync"
)

type contextTag int32

const (
	waitGroup contextTag = iota
)

func WithWaitGroup(ctx context.Context, wg *sync.WaitGroup) context.Context {
	return context.WithValue(ctx, waitGroup, wg)
}

func IncAt(ctx context.Context) bool {
	wg, ok := ctx.Value(waitGroup).(*sync.WaitGroup)
	if !ok {
		return ok
	}
	wg.Add(1)
	return true
}

func DoneFrom(ctx context.Context) bool {
	wg, ok := ctx.Value(waitGroup).(*sync.WaitGroup)
	if !ok {
		return ok
	}
	wg.Done()
	return true
}
