package config

import (
	"strings"

	"github.com/mediocregopher/radix/v3"
	"github.com/sirupsen/logrus"
)

type RedisConfigStore struct {
	Pool *radix.Pool
}

func (rs *RedisConfigStore) GetValue(key string) interface{} {
	prefixStripped := strings.TrimPrefix(key, "sgpdb.")

	var v string
	err := rs.Pool.Do(radix.Cmd(&v, "HGET", "sgpdb_config", prefixStripped))
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
	prefixStripped := strings.TrimPrefix(key, "sgpdb.")

	err := rs.Pool.Do(radix.Cmd(nil, "HSET", "sgpdb_config", prefixStripped, value))
	if err != nil {
		return err
	}

	return nil
}

func (e *RedisConfigStore) Name() string {
	return "redis"
}
