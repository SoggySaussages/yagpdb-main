package common

import (
	"encoding/json"

	"github.com/botlabs-gg/yagpdb/v2/common/redis"
)

// GetRedisJson executes a get redis command and unmarshals the value into out
func GetRedisJson(key string, out interface{}) error {
	var resp []byte
	err := RedisPool.Do(redis.Cmd(&resp, "GET", key))
	if err != nil {
		return err
	}

	if len(resp) == 0 {
		return nil
	}

	err = json.Unmarshal(resp, out)
	return err
}

// SetRedisJson marshals the value and runs a set redis command for key
func SetRedisJson(key string, value interface{}) error {
	serialized, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = RedisPool.Do(redis.Cmd(nil, "SET", key, string(serialized)))
	return err
}

func MultipleCmds(cmds ...redis.RedisCmdAction) error {
	for _, v := range cmds {
		err := RedisPool.Do(v)
		if err != nil {
			return err
		}
	}

	return nil
}
