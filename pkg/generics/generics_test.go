package generics

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type Resource struct {
	ID uuid.UUID
}

type Deps struct {
	Res *Resource `groat:"resource"`
}

func TestInjector(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		exp := &Resource{ID: uuid.New()}

		res := Injector[*Resource, Deps](t, exp, Deps{}, "resource")
		assert.Equal(t, exp, res.Res)
	})
	t.Run("can't inject", func(t *testing.T) {
		exp := &Resource{ID: uuid.New()}

		res := Injector[*Resource, Deps](t, exp, Deps{}, "another")
		assert.Nil(t, res.Res)
	})
}
