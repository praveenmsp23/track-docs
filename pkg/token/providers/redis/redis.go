package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/praveenmsp23/trackdocs/pkg/cache"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/logger"
	"github.com/praveenmsp23/trackdocs/pkg/token/base"
)

const (
	TokenPrefix = "token_v1::"
)

func getTokenCacheKey(tid string) string {
	return fmt.Sprintf("%s%s", TokenPrefix, tid)
}

func GetProvider(cfg *config.Config) (*RedisProvider, error) {
	c, err := cache.NewCache(cfg)
	if err != nil {
		return nil, err
	}
	return &RedisProvider{client: c, maxlifetime: cfg.TokenLifeTime}, nil
}

type RedisTokenStore struct {
	tid         string
	client      *cache.Cache
	maxlifetime int64
}

func (st *RedisTokenStore) Set(key, value string) error {
	return st.client.HSetString(getTokenCacheKey(st.tid), key, value)
}

func (st *RedisTokenStore) Get(key string) (string, bool) {
	res, err := st.client.HGetString(getTokenCacheKey(st.tid), key)
	return res, err == nil
}

func (st *RedisTokenStore) GetAll() (map[string]string, error) {
	return st.client.HGetAll(getTokenCacheKey(st.tid))
}

func (st *RedisTokenStore) Delete(key string) error {
	return st.client.HDel(getTokenCacheKey(st.tid), key)
}

func (st *RedisTokenStore) TokenID() string {
	return st.tid
}

type RedisProvider struct {
	client      *cache.Cache
	maxlifetime int64
}

func (pder *RedisProvider) TokenInit(tid string) (base.Token, error) {
	key := getTokenCacheKey(tid)
	logger.Debugf("TokenInit key:%s", key)
	logger.Debugf("TokenInit maxlifetime:%d", pder.maxlifetime)
	err := pder.client.HSetString(key, "active", "1")
	if err != nil {
		return nil, err
	}
	_, err = pder.client.Expire(key, time.Duration(time.Second*time.Duration(pder.maxlifetime)))
	if err != nil {
		return nil, err
	}
	return &RedisTokenStore{tid: tid, maxlifetime: pder.maxlifetime, client: pder.client}, nil
}

func (pder *RedisProvider) TokenRead(tid string) (base.Token, error) {
	t, err := pder.client.HGetString(getTokenCacheKey(tid), "active")
	if err != nil {
		return nil, err
	}
	if t != "1" {
		return nil, errors.New("token not available")
	}
	return &RedisTokenStore{tid: tid, maxlifetime: pder.maxlifetime, client: pder.client}, nil
}

func (pder *RedisProvider) TokenDestroy(tid string) error {
	err := pder.client.Del(getTokenCacheKey(tid))
	if err != nil {
		return err
	}
	return nil
}

func (pder *RedisProvider) TokenGC(maxlifetime int64) {
}
