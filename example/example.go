package example

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

type Config struct {
	TTL      time.Duration
	Capacity int
}

type CacheService[K comparable, V any] struct {
	lock        *sync.RWMutex
	data        map[K]V
	ctx         context.Context
	ticker      *time.Ticker
	cfg         Config
	remoteCache DB[K, V]
}

//go:generate mockery --name DB
type DB[K comparable, V any] interface {
	Set(k K, v V) error
	Get(k K) (v V, ok error)
}

func NewCacheService[K comparable, V any](ctx context.Context, cfg Config, db DB[K, V]) *CacheService[K, V] {
	srv := CacheService[K, V]{
		lock:        new(sync.RWMutex),
		ctx:         ctx,
		cfg:         cfg,
		ticker:      time.NewTicker(cfg.TTL),
		data:        make(map[K]V, cfg.Capacity),
		remoteCache: db,
	}
	go srv.loop()

	return &srv
}

func (srv *CacheService[K, V]) Set(k K, v V) error {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	if err := srv.remoteCache.Set(k, v); err != nil {
		return err
	}
	srv.data[k] = v
	return nil
}

func (srv *CacheService[K, V]) localGet(k K) (v V, ok bool) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	v, ok = srv.data[k]
	return v, ok
}

func (srv *CacheService[K, V]) Get(k K) (v V, err error) {
	v, ok := srv.localGet(k)
	if ok {
		return v, nil
	}

	srv.lock.Lock()
	defer srv.lock.Unlock()
	v, err = srv.remoteCache.Get(k)
	if err != nil {
		return v, err
	}
	srv.data[k] = v
	return v, nil
}

func (srv *CacheService[K, V]) loop() {
	for {
		select {
		case <-srv.ctx.Done():
			return
		case <-srv.ticker.C:
			srv.beat()
			srv.ticker.Reset(srv.cfg.TTL)
		}
	}
}

func (srv *CacheService[K, V]) beat() {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	srv.data = make(map[K]V, srv.cfg.Capacity)
}
