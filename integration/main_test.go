package integration

import (
	"context"
	"os"
	"testing"

	"github.com/godepo/groat"
	"github.com/godepo/groat/pkg/ctxgroup"
	"github.com/godepo/groat/pkg/generics"
	"github.com/stretchr/testify/require"
)

type Resource struct {
	ID string
}

type Deps struct {
	Res *Resource `groat:"test_resource"`
}

type State struct {
}

type SUT struct {
	Depend *Resource
}

var suite *Container[Deps, State, *SUT]

func TestMain(m *testing.M) {
	suite = New[Deps, State, *SUT](m, func(t *testing.T) *groat.Case[Deps, State, *SUT] {
		tcs := groat.New[Deps, State, *SUT](t, func(t *testing.T, deps Deps) *SUT {
			return &SUT{Depend: deps.Res}
		})
		return tcs
	}, func(ctx context.Context) (Injector[Deps], error) {
		ctxgroup.IncAt(ctx)
		go func() {
			defer ctxgroup.DoneFrom(ctx)
			select {
			case <-ctx.Done():
				return
			}
		}()
		return func(t *testing.T, injectTo Deps) Deps {
			res := &Resource{ID: "OK"}
			return generics.Injector[*Resource, Deps](t, res, injectTo, "test_resource")
		}, nil
	})
	os.Exit(suite.Go())
}

func TestIntegration(t *testing.T) {
	t.Run("should be injected resource", func(t *testing.T) {
		tcs := suite.Case(t)
		require.NotNil(t, tcs.SUT)
		require.NotNil(t, tcs.Deps.Res)
		require.NotNil(t, tcs.SUT.Depend)
		require.Equal(t, tcs.Deps.Res, tcs.SUT.Depend)
		require.Equal(t, tcs.Deps.Res.ID, tcs.SUT.Depend.ID)
		require.Equal(t, tcs.SUT.Depend.ID, "OK")
	})
}
