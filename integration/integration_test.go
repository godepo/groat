package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/godepo/groat"
	"github.com/stretchr/testify/assert"
)

type ContainerState struct {
	ExpectError error
}

type ContainerSut = *Container[Deps, State, *SUT]

func newContainerCase(t *testing.T) *groat.Case[Deps, ContainerState, ContainerSut] {
	tcs := groat.New[Deps, ContainerState, ContainerSut](t, func(t *testing.T, deps Deps) ContainerSut {
		sut := &Container[Deps, State, *SUT]{}
		return sut
	})
	tcs.Go()
	return tcs
}

func TestContainer_Go(t *testing.T) {
	t.Run("should be able to be failed", func(t *testing.T) {
		t.Run("failed when bootstrapper return error", func(t *testing.T) {
			tcs := newContainerCase(t)
			tcs.Given(
				ArrangeBootstrapperWithError(tcs.SUT),
			)
			assert.Equal(t, 1, tcs.SUT.Go())
		})

		t.Run("failed when bootstrapper failed in panic", func(t *testing.T) {
			tcs := newContainerCase(t)
			tcs.Given(
				ArrangeBootstrapperWithPanic(tcs.SUT),
			)
			assert.Equal(t, 100, tcs.SUT.Go())
		})
	})
}

func ArrangeBootstrapperWithError(sut ContainerSut) groat.Given[ContainerState] {
	return func(t *testing.T, state ContainerState) ContainerState {
		sut.bootstraps = append(sut.bootstraps, func(ctx context.Context) (Injector[Deps], error) {
			return nil, errors.New("failed bootstrapper")
		})
		return state
	}
}

func ArrangeBootstrapperWithPanic(sut ContainerSut) groat.Given[ContainerState] {
	return func(t *testing.T, state ContainerState) ContainerState {
		sut.bootstraps = append(sut.bootstraps, func(ctx context.Context) (Injector[Deps], error) {
			panic(errors.New("panic bootstrapper"))
		})
		return state
	}
}
