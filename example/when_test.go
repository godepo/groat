package example

import "testing"

func ActGetKey(t *testing.T, deps CacheDeps, state CacheState) CacheState {
	t.Helper()
	deps.Remote.EXPECT().Get(state.Key).Return(state.Value, nil)
	return state
}

func ActSetKeyToRemoteCache(t *testing.T, deps CacheDeps, state CacheState) CacheState {
	t.Helper()
	deps.Remote.EXPECT().Set(state.Key, state.Value).Return(nil)
	return state
}

func ActSetKeyToRemoteCacheFailByExpectedError(_ *testing.T, deps CacheDeps, state CacheState) CacheState {
	deps.Remote.EXPECT().Set(state.Key, state.Value).Return(state.ExpectError)
	return state
}

func ActGetKeyFailedAtRemoteCache(_ *testing.T, deps CacheDeps, state CacheState) CacheState {
	deps.Remote.EXPECT().Get(state.Key).Return(state.Value, state.ExpectError)
	return state
}
