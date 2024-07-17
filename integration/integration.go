package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/godepo/groat"
	"github.com/godepo/groat/pkg/ctxgroup"
)

// Injector is type of function for inject dependencies to target struct.
type Injector[T any] func(t *testing.T, injectTo T) T

// Bootstrap prepare injectors for test cases.
type Bootstrap[T any] func(ctx context.Context) (Injector[T], error)

// Provider construct groat.Case for each running test.
type Provider[Deps any, State any, SUT any] func(t *testing.T) *groat.Case[Deps, State, SUT]

type Container[Deps any, State any, SUT any] struct {
	injectors  []Injector[Deps]
	bootstraps []Bootstrap[Deps]
	provider   Provider[Deps, State, SUT]
	m          *testing.M
}

func New[Deps any, State any, SUT any](
	m *testing.M,
	provider Provider[Deps, State, SUT],
	all ...Bootstrap[Deps],
) *Container[Deps, State, SUT] {
	return &Container[Deps, State, SUT]{
		m:          m,
		provider:   provider,
		bootstraps: all,
	}
}

func (c *Container[Deps, State, SUT]) Go() (res int) {
	defer func() {
		rr := recover()
		if rr != nil {
			fmt.Printf("unexpected panic: %v\n", rr) //nolint:forbidigo
			res = 100
		}
	}()
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = ctxgroup.WithWaitGroup(ctx, wg)
	defer func() {
		cancel()
		wg.Wait()
	}()

	for i, b := range c.bootstraps {
		injector, err := b(ctx)
		if err != nil {
			fmt.Printf("=== GROAT fail bootstrap[%d]: %v\n", i, err) //nolint:forbidigo
			return 1
		}
		c.injectors = append(c.injectors, injector)
	}

	return c.m.Run()
}

func (c *Container[Deps, State, SUT]) Case(t *testing.T) *groat.Case[Deps, State, SUT] {
	t.Helper()
	tcs := c.provider(t)
	before := make([]groat.Before[Deps], 0, len(c.injectors))
	for _, injector := range c.injectors {
		before = append(before, func(t *testing.T, deps Deps) Deps {
			t.Helper()
			deps = injector(tcs.T, deps)
			return deps
		})
	}
	tcs.Before(before...)
	tcs.Go()
	return tcs
}
