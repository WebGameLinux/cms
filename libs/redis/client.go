package redis

import (
		"github.com/WebGameLinux/cms/configuration"
		"github.com/go-redis/redis"
)

func GetInstance(name ...string) *redis.Client {
		if len(name) == 0 {
				name = append(name, TempDb)
		}
		var instance = newInstance(name[0])
		return instance
}

func newInstance(name string) *redis.Client {
		instance := redis.NewClient(getOptions(name))
		return instance
}

func getOptions(name string) *redis.Options {
		kv := configuration.GetRedisKvCnf()
		return &redis.Options{
				Addr:     kv.Conn(),
				Password: kv.Password(),
				DB:       getDbByName(name),
		}
}
