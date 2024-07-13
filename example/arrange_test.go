package example

import (
	"errors"
	"testing"

	"github.com/godepo/groat"
)

func ArrangeRandomKeyValue(_ *testing.T, state CacheState) CacheState {
	state.Key = state.Faker.Music().Name()
	state.Value = state.Faker.Int()
	return state
}

func ArrangeExpectedError(_ *testing.T, state CacheState) CacheState {
	state.ExpectError = errors.New(state.Faker.UUID().V4())
	return state
}

func PopulateKeyToCache(sut *CacheService[string, int]) groat.Given[CacheState] {
	return func(t *testing.T, state CacheState) CacheState {
		sut.lock.Lock()
		defer sut.lock.Unlock()

		sut.data[state.Key] = state.Value
		return state
	}
}
