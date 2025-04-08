package config

import (
	"strings"

	"github.com/botlabs-gg/yagpdb/v2/common/redis"
	"github.com/sirupsen/logrus"
)

type RedisConfigStore struct {
	Pool *redis.RedisPool
}

func (rs *RedisConfigStore) GetValue(key string) interface{} {
	prefixStripped := strings.TrimPrefix(key, "yagpdb.")

	var v string
	err := rs.Pool.Do(redis.Cmd(&v, "HGET", "yagpdb_config", prefixStripped))
	if err != nil {
		logrus.WithError(err).Error("[redis_config_source] failed retrieving value")
		return nil
	}

	if v == "" {
		return nil
	}

	return v
}

func (rs *RedisConfigStore) SaveValue(key, value string) error {
	prefixStripped := strings.TrimPrefix(key, "yagpdb.")

	err := rs.Pool.Do(redis.Cmd(nil, "HSET", "yagpdb_config", prefixStripped, value))
	if err != nil {
		return err
	}

	return nil
}

func (e *RedisConfigStore) Name() string {
	return "redis"
}
