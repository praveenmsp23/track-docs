package lock

import (
	"time"

	"github.com/praveenmsp23/trackdocs/pkg/cache"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/crypto"
)

type RedisLock struct {
	cache *cache.Cache
	cfg   *config.Config
}

func NewRedisLock(cache *cache.Cache, cfg *config.Config) (*RedisLock, error) {
	return &RedisLock{cache: cache, cfg: cfg}, nil
}

func (l *RedisLock) NewMutex(name string, options ...Option) *Mutex {
	m := &Mutex{
		name:       name,
		expiry:     10 * time.Second,
		retryCount: 50,
		retryDelay: 1 * time.Second,
		cache:      l.cache,
		cfg:        l.cfg,
	}
	for _, o := range options {
		o.Apply(m)
	}
	if m.value == "" {
		m.value = crypto.GenerateId("lok", 32)
	}
	return m
}

// An Option configures a mutex.
type Option interface {
	Apply(*Mutex)
}

// OptionFunc is a function that configures a mutex.
type OptionFunc func(*Mutex)

// Apply calls f(mutex)
func (f OptionFunc) Apply(mutex *Mutex) {
	f(mutex)
}

// WithExpiry can be used to set the expiry of a mutex to the given value.
func WithExpiry(expiry time.Duration) Option {
	return OptionFunc(func(m *Mutex) {
		m.expiry = expiry
	})
}

// WithRetryCount can be used to set the number of times lock acquire is attempted.
func WithRetryCount(retryCount int) Option {
	return OptionFunc(func(m *Mutex) {
		m.retryCount = retryCount
	})
}

// WithRetryDelay can be used to set the amount of time to wait between retries.
func WithRetryDelay(delay time.Duration) Option {
	return OptionFunc(func(m *Mutex) {
		m.retryDelay = delay
	})
}

// WithValue can be used to assign the random value without having to call lock.
// This allows the ownership of a lock to be "transferred" and allows the lock to be unlocked from elsewhere.
func WithValue(v string) Option {
	return OptionFunc(func(m *Mutex) {
		m.value = v
	})
}
