package services

import (
		"github.com/WebGameLinux/cms/configuration"
		"github.com/WebGameLinux/cms/libs/redis"
		"github.com/astaxie/beego/cache"
		redis2 "github.com/go-redis/redis"
		"time"
)

type RedisServiceInterface interface {
		cache.Cache
}

type RedisService struct {
		client *redis2.Client
}

func (this *RedisService) SetClient(client *redis2.Client) RedisServiceInterface {
		this.client = client
		return this
}

func (this *RedisService) Get(key string) interface{} {
		if str, err := this.client.Get(key).Result(); err == nil {
				return str
		}
		return nil
}

func (this *RedisService) GetMulti(keys []string) []interface{} {
		if v, err := this.client.MGet(keys...).Result(); err == nil {
				return v
		}
		return nil
}

func (this *RedisService) Put(key string, val interface{}, timeout time.Duration) error {
		if ok, err := this.client.SetXX(key, val, timeout).Result(); err != nil || !ok {
				return err
		}
		return nil
}

func (this *RedisService) Delete(key string) error {
		if _, err := this.client.Del(key).Result(); err != nil {
				return err
		}
		return nil
}

func (this *RedisService) Incr(key string) error {
		if _, err := this.client.Incr(key).Result(); err != nil {
				return err
		}
		return nil
}

func (this *RedisService) Keep(key string, expire time.Duration) error {
		if _, err := this.client.Expire(key, expire).Result(); err != nil {
				return err
		}
		return nil
}

func (this *RedisService) SetEx(key string, value interface{}, timeout time.Duration) bool {
		if v := this.client.SetNX(key, value, timeout); v != nil {
				ok, err := v.Result()
				if err != nil || !ok {
						return false
				}
		}
		return false
}

func (this *RedisService) Decr(key string) error {
		if _, err := this.client.Decr(key).Result(); err != nil {
				return err
		}
		return nil
}

func (this *RedisService) IsExist(key string) bool {
		if _, err := this.client.Exists(key).Result(); err != nil {
				return false
		}
		return true
}

func (this *RedisService) ClearAll() error {
		cachedKeys, err := this.client.Keys("*").Result()
		if err != nil {
				return err
		}
		if _, err = this.client.Del(cachedKeys...).Result(); err != nil {
				return err
		}
		return nil
}

func (this *RedisService) StartAndGC(config string) error {
		if this.client == nil {
				this.client = redis.GetInstance()
		}
		return nil
}

func (this *RedisService) GetRedis(name ...string) *redis2.Client {
		if len(name) == 0 {
				return this.client
		}
		return redis.GetInstance(name[0])
}

func (this *RedisService) Ok(state string) bool {
		return state == "OK"
}

func NewRedisService() *RedisService {
		service := new(RedisService)
		_ = service.StartAndGC(configuration.GetRedisKvCnf().String())
		return service
}

func CreateRedisService(name string) *RedisService {
		service := new(RedisService)
		service.client = redis.GetInstance(name)
		return service
}
