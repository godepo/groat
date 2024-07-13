package example

import (
	"testing"
	"time"

	"github.com/godepo/groat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func WaitLongTime(duration time.Duration) groat.Then[CacheState] {
	return func(t *testing.T, state CacheState) {
		time.Sleep(duration)
	}
}

func ExpectNoError(t *testing.T, state CacheState) {
	require.NoError(t, state.ResultError)
}

func ExpectEmptyKey(sut *CacheService[string, int]) groat.Then[CacheState] {
	return func(t *testing.T, state CacheState) {
		sut.lock.Lock()
		defer sut.lock.Unlock()
		_, ok := sut.data[state.Key]
		require.False(t, ok)
	}
}

func AssertKeyValueAsGiven(sut *CacheService[string, int]) groat.Then[CacheState] {
	return func(t *testing.T, state CacheState) {
		sut.lock.Lock()
		defer sut.lock.Unlock()
		result, ok := sut.data[state.Key]
		require.True(t, ok)
		assert.Equal(t, state.Value, result)
	}
}

func ExpectResultError(t *testing.T, state CacheState) {
	require.ErrorIs(t, state.ExpectError, state.ResultError)
}
