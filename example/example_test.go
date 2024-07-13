package example

import (
	"context"
	"testing"
	"time"

	"github.com/godepo/groat"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/require"
)

type CacheDeps struct {
	cfg    Config
	ctx    context.Context
	cancel context.CancelFunc
	Faker  faker.Faker
	Remote *MockDB[string, int]
}

type CacheState struct {
	Key         string
	Faker       faker.Faker
	Value       int
	Result      int
	ResultError error
	ExpectError error
}

type testCase = *groat.Case[CacheDeps, CacheState, *CacheService[string, int]]

func newTestCase(t *testing.T) testCase {
	fxt := groat.New[CacheDeps, CacheState, *CacheService[string, int]](
		t,
		func(t *testing.T, deps CacheDeps) *CacheService[string, int] {
			return NewCacheService[string, int](deps.ctx, deps.cfg, deps.Remote)
		},
		func(t *testing.T, deps CacheDeps) CacheDeps {
			deps.Remote = NewMockDB[string, int](t)

			deps.cfg = Config{
				TTL:      time.Millisecond * 25,
				Capacity: 1000,
			}
			deps.ctx = context.Background()
			deps.ctx, deps.cancel = context.WithCancel(deps.ctx)
			deps.Faker = faker.New()
			return deps
		},
	)
	fxt.Given(func(t *testing.T, state CacheState) CacheState {
		state.Faker = fxt.Deps.Faker
		return state
	})
	fxt.After(func(t *testing.T, deps CacheDeps) {
		deps.cancel()
	})
	fxt.Go()
	return fxt
}

func TestCacheService_Set(t *testing.T) {
	t.Run("should be able to add and retrieve set item", func(t *testing.T) {
		fxt := newTestCase(t)

		fxt.Given(ArrangeRandomKeyValue).
			When(ActSetKeyToRemoteCache).
			Then(AssertKeyValueAsGiven(fxt.SUT))

		require.NoError(t, fxt.SUT.Set(fxt.State.Key, fxt.State.Value))
	})
	t.Run("should be able not retrieve item", func(t *testing.T) {
		fxt := newTestCase(t)
		fxt.
			Given(ArrangeRandomKeyValue).
			When(ActSetKeyToRemoteCache).
			Then(
				WaitLongTime(time.Millisecond*100),
				ExpectEmptyKey(fxt.SUT),
				ExpectNoError,
			)

		fxt.State.ResultError = fxt.SUT.Set(fxt.State.Key, fxt.State.Value)
	})
	t.Run("should be able to failed", func(t *testing.T) {
		t.Run("when remote service failed", func(t *testing.T) {
			fxt := newTestCase(t)

			fxt.
				Given(ArrangeRandomKeyValue, ArrangeExpectedError).
				When(ActSetKeyToRemoteCacheFailByExpectedError).
				Then(
					ExpectEmptyKey(fxt.SUT),
					ExpectResultError,
				)
			fxt.State.ResultError = fxt.SUT.Set(fxt.State.Key, fxt.State.Value)
		})
	})
}

func TestCacheService_Get(t *testing.T) {
	t.Run("should be able to no retrieve item from remote cache, when local is empty", func(t *testing.T) {
		fxt := newTestCase(t)
		fxt.
			Given(ArrangeRandomKeyValue).
			When(ActGetKey).
			Then(AssertKeyValueAsGiven(fxt.SUT))

		fxt.State.Result, fxt.State.ResultError = fxt.SUT.Get(fxt.State.Key)
	})

	t.Run("should be able to retrieve item from local cache", func(t *testing.T) {
		fxt := newTestCase(t)
		fxt.
			Given(ArrangeRandomKeyValue, PopulateKeyToCache(fxt.SUT)).
			Then(AssertKeyValueAsGiven(fxt.SUT))

		fxt.State.Result, fxt.State.ResultError = fxt.SUT.Get(fxt.State.Key)
	})

	t.Run("should be able to be failed", func(t *testing.T) {
		t.Run("when remote service failed", func(t *testing.T) {
			fxt := newTestCase(t)
			fxt.
				Given(ArrangeRandomKeyValue, ArrangeExpectedError).
				When(ActGetKeyFailedAtRemoteCache).
				Then(ExpectResultError, ExpectEmptyKey(fxt.SUT))

			fxt.State.Result, fxt.State.ResultError = fxt.SUT.Get(fxt.State.Key)
		})
	})
}

func TestShutdown(t *testing.T) {
	t.Run("should be able to shutdown", func(t *testing.T) {
		fxt := newTestCase(t)
		fxt.
			Given(ArrangeRandomKeyValue).
			After(func(t *testing.T, deps CacheDeps) {
				t.Helper()
				fxt.SUT.lock.Lock()
				fxt.SUT.data[fxt.State.Key] = fxt.State.Value
				fxt.SUT.lock.Unlock()
				time.Sleep(time.Millisecond * 100)
				AssertKeyValueAsGiven(fxt.SUT)(t, fxt.State)
			})
	})
}
