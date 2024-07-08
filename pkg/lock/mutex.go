package lock

// Refer https://redis.io/docs/reference/patterns/distributed-locks/

import (
	"errors"
	"time"

	"github.com/praveenmsp23/trackdocs/pkg/cache"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/crypto"
)

const (
	unlockScript = `
		if redis.call("get",KEYS[1]) == ARGV[1] then
		    return redis.call("del",KEYS[1])
		else
		    return 0
		end
	`
)

var ErrUnLockFailed = errors.New("lock unlock failed")

// A Mutex is a distributed mutual exclusion lock.
type Mutex struct {
	name string

	retryCount int
	retryDelay time.Duration

	expiry time.Duration
	cache  *cache.Cache
	cfg    *config.Config

	value string
}

// Name returns mutex name
func (m *Mutex) Name() string {
	return m.name
}

// Lock attempts to put a lock on the key for a mutex expiry duration.
// If the lock was successfully acquired, true will be returned.
func (m *Mutex) Lock() (bool, error) {
	for i := 0; i < m.retryCount; i++ {
		success, err := m.acquire()
		if success || err != nil {
			return success, err
		}
		if i == m.retryCount-1 {
			return false, nil
		}
		// Wait a random delay before to retry
		time.Sleep(crypto.GenerateRandomDuration(m.retryDelay))
	}
	return false, nil
}

// Unlock attempts to remove the lock on a key so long as the value matches.
// If the lock cannot be removed, either because the key has already expired or
// because the value was incorrect, an error will be returned.
func (m *Mutex) Unlock() error {
	return m.release(m.name)
}

func (m *Mutex) acquire() (bool, error) {
	reply, err := m.cache.SetNX(m.name, m.value, m.expiry)
	if err != nil {
		return false, err
	}
	return reply, nil
}

func (m *Mutex) release(lockName string) error {
	res := m.cache.Eval(unlockScript, []string{lockName}, m.value)
	if res.Err() != nil {
		return res.Err()
	}
	result, err := res.Int()
	if err != nil {
		return err
	}
	if result == 0 {
		return ErrUnLockFailed
	}
	return nil
}
