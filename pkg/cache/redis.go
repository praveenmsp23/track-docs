package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"

	"encoding/json"

	"github.com/praveenmsp23/trackdocs/pkg/config"
)

type Cache struct {
	client  *redis.Client
	limiter *Limiter
	ctx     context.Context
}

const (
	DefaultConnectionTimeout = 2 * time.Minute
	DefaultExpiry            = 604800
)

func NewCache(cfg *config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.CacheSource,
		Password: cfg.CacheSourcePassword,
		DB:       0,
	})
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	limiter := NewLimiter(client)
	return &Cache{client: client, ctx: ctx, limiter: limiter}, nil
}

// GetClient get internal redis client
func (c *Cache) GetClient() *redis.Client {
	return c.client
}

// Allow is a shortcut for AllowN(ctx, key, limit, 1).
func (c *Cache) Allow(key string, limit int) (*Result, error) {
	return c.limiter.Allow(c.ctx, key, PerMinute(limit))
}

// AllowN reports whether n events may happen at time now.
func (c *Cache) AllowN(key string, limit, n int) (*Result, error) {
	return c.limiter.AllowN(c.ctx, key, PerMinute(limit), n)
}

// Set sets the cache value for the key
func (c *Cache) Set(key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.client.Set(c.ctx, key, string(val), DefaultExpiry*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetX sets the cache value for the key
func (c *Cache) SetX(key string, value interface{}, expiry time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.client.Set(c.ctx, key, string(val), expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

// Del deletes the cache value for the key
func (c *Cache) Del(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

// HDel deletes the cache value for the key
func (c *Cache) HDel(key, field string) error {
	err := c.client.HDel(c.ctx, key, field).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get get the cache value for the key
func (c *Cache) Get(key string, out interface{}) error {
	res := c.client.Get(c.ctx, key)

	if res.Err() != nil {
		return res.Err()
	}
	bytes, err := res.Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}
	return err
}

func (c *Cache) Exists(key string) (bool, error) {
	res := c.client.Exists(c.ctx, key)

	if res.Err() != nil {
		return false, res.Err()
	}
	val, err := res.Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (c *Cache) setAny(key string, value interface{}) error {
	res := c.client.Set(c.ctx, key, value, DefaultExpiry*time.Second)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (c *Cache) setAnyX(key string, value interface{}, expiry time.Duration) error {
	res := c.client.Set(c.ctx, key, value, expiry)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

// SetString sets the cache value for the key
func (c *Cache) SetString(key, value string) error {
	return c.setAny(key, value)
}

// SetInt sets the cache value for the key
func (c *Cache) SetInt(key string, value int) error {
	return c.setAny(key, value)
}

// SetInt64 sets the cache value for the key
func (c *Cache) SetInt64(key string, value int64) error {
	return c.setAny(key, value)
}

// SetXString sets the cache value for the key
func (c *Cache) SetXString(key, value string, expiry time.Duration) error {
	return c.setAnyX(key, value, expiry)
}

// SetXInt sets the cache value for the key
func (c *Cache) SetXInt(key string, value int, expiry time.Duration) error {
	return c.setAnyX(key, value, expiry)
}

// SetXInt64 sets the cache value for the key
func (c *Cache) SetXInt64(key string, value int64, expiry time.Duration) error {
	return c.setAnyX(key, value, expiry)
}

// SetNX Redis `SET key value [expiration] NX` command.
func (c *Cache) SetNX(key, value string, expiry time.Duration) (bool, error) {
	res := c.client.SetNX(c.ctx, key, value, expiry)
	if res.Err() != nil {
		return false, res.Err()
	}
	return res.Val(), nil
}

// Expire Redis `EXPIRE key [expiration]` command.
func (c *Cache) Expire(key string, expiry time.Duration) (bool, error) {
	res := c.client.Expire(c.ctx, key, expiry)
	if res.Err() != nil {
		return false, res.Err()
	}
	return res.Val(), nil
}

// PTTL Redis `PTTL key` command.
func (c *Cache) PTTL(key string) (time.Duration, error) {
	res := c.client.PTTL(c.ctx, key)
	if res.Err() != nil {
		return time.Millisecond, res.Err()
	}
	return res.Val(), nil
}

// HGetString get the cache value for the key
func (c *Cache) HGetString(key, field string) (string, error) {
	res := c.client.HGet(c.ctx, key, field)
	if res.Err() != nil {
		return "", res.Err()
	}
	out, err := res.Result()
	if err != nil {
		return "", err
	}
	return out, nil
}

// HGetAll
func (c *Cache) HGetAll(key string) (map[string]string, error) {
	res := c.client.HGetAll(c.ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	out, err := res.Result()
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HGetInt get the cache value for the key
func (c *Cache) HGetInt(key, field string) (int, error) {
	res := c.client.HGet(c.ctx, key, field)
	if res.Err() != nil {
		return 0, res.Err()
	}
	out, err := res.Int()
	if err != nil {
		return 0, err
	}
	return out, nil
}

// HGetInt64 get the cache value for the key
func (c *Cache) HGetInt64(key, field string) (int64, error) {
	res := c.client.HGet(c.ctx, key, field)
	if res.Err() != nil {
		return 0, res.Err()
	}
	out, err := res.Int64()
	if err != nil {
		return 0, err
	}
	return out, nil
}

// GetString get the cache value for the key
func (c *Cache) GetString(key string) (string, error) {
	res := c.client.Get(c.ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}
	out, err := res.Result()
	if err != nil {
		return "", err
	}
	return out, nil
}

// GetInt get the cache value for the key
func (c *Cache) GetInt(key string) (int, error) {
	res := c.client.Get(c.ctx, key)
	if res.Err() != nil {
		return 0, res.Err()
	}
	out, err := res.Int()
	if err != nil {
		return 0, err
	}
	return out, nil
}

// GetInt64 get the cache value for the key
func (c *Cache) GetInt64(key string) (int64, error) {
	res := c.client.Get(c.ctx, key)
	if res.Err() != nil {
		return 0, res.Err()
	}
	out, err := res.Int64()
	if err != nil {
		return 0, err
	}
	return out, nil
}

// FlushAll flushes all cache. **WARNING** only for development
func (c *Cache) FlushAll() error {
	res := c.client.FlushAll(c.ctx)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

// HSet sets the cache value for the key in hash
func (c *Cache) HSet(key, field string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.client.HSet(c.ctx, key, field, string(val)).Err()
	if err != nil {
		return err
	}
	return nil
}

// HSet sets the cache value for the key in hash
func (c *Cache) HSetString(key, field, value string) error {
	err := c.client.HSet(c.ctx, key, field, value).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get get the cache value for the key
func (c *Cache) HGet(key, field string, out interface{}) error {
	res := c.client.HGet(c.ctx, key, field)

	if res.Err() != nil {
		return res.Err()
	}
	bytes, err := res.Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}
	return err
}

// Keys returns the keys matching pattern
func (c *Cache) Keys(pattern string) ([]string, error) {
	res := c.client.Keys(c.ctx, pattern)
	if res.Err() != nil {
		return []string{}, res.Err()
	}
	return res.Val(), nil
}

// Eval evaluates the lua script
func (c *Cache) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	return c.client.Eval(c.ctx, script, keys, args...)
}
