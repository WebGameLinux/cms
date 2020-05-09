package configuration

import (
		"fmt"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/astaxie/beego/config"
)

type RedisConfiguration struct {
		CacheConfiguration
}

type RedisConfigurationWrapper interface {
		Port() int
		Password() string
		User() string
		Host() string
		Collection() string
		Conn() string
		DbNum() int

		Key(string) string
}

var redisKvCnf *RedisConfiguration

const (
		CnfKvRedisDbNum             = "db"
		CnfKvRedisHost              = "host"
		CnfKvRedisPort              = "port"
		CnfKvRedisPassword          = "password"
		CnfKvRedisUser              = "user"
		CnfKvRedisCollection        = "collection"
		CnfKvRedisScopeDefault      = "default"
		CnfKvRedisCollectionDefault = "default"
		CnfKvRedisGlobal            = "cache.redis"
		CnfKvRedisHostDefault       = "127.0.0.1"
		CnfKvRedisPortDefault       = 6379
		CnfKvRedisDbNumDefault      = 0
		StringRedisTemplate         = `{"key":"%s","conn":"%s","dbNum":"%d","user":"%s","password":"%s"}`
)

func RedisCnfString() string {
		return redisKvCnf.String()
}

func RedisCnfScope(name string) CacheConfigure {
		return redisKvCnf.Scope(name)
}

func GetRedisKvCnf(Properties ...interface{}) *RedisConfiguration {
		if redisKvCnf == nil {
				redisKvCnf = new(RedisConfiguration)
				redisKvCnf.init()
		}
		if len(Properties) > 0 {
				kv := Properties[0]
				switch kv.(type) {
				case mapper.Mapper:
						if m, ok := kv.(mapper.Mapper); ok {
								redisKvCnf.KvCnf = m
						}
				case config.Configer:
						if m, ok := kv.(config.Configer); ok {
								redisKvCnf.KvCnf = AppConfig2Map(m, CnfKvRedisGlobal)
						}
				default:
						redisKvCnf.KvCnf = map[string]interface{}{}
				}
		}
		return redisKvCnf
}

func (this *RedisConfiguration) init() {
		if this.KvCnf == nil {
				this.KvCnf = make(map[string]interface{})
		}
		if this.ScopeName == "" {
				this.ScopeName = CnfKvRedisGlobal
		}
}

func (this *RedisConfiguration) String() string {
		return fmt.Sprintf(StringRedisTemplate, this.Args()...)
}

func (this *RedisConfiguration) Args() []interface{} {
		return []interface{}{this.Collection(),
				this.Conn(),
				this.DbNum(),
				this.User(),
				this.Password(),
		}
}

func (this *RedisConfiguration) Scope(name string) CacheConfigure {
		kv := new(RedisConfiguration)
		kv.init()
		if name != "" {
				kv.ScopeName = this.ScopeName + "." + name
		}
		kv.KvCnf = this.KvCnf
		return kv
}

func (this *RedisConfiguration) Get(key string, def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		kv := mapper.Mapper(this.KvCnf)
		return kv.Get(this.Key(key), def...)
}

func (this *RedisConfiguration) User() string {
		return this.Get(CnfKvRedisUser, "")
}

func (this *RedisConfiguration) Password() string {
		return this.Get(CnfKvRedisPassword, "")
}

func (this *RedisConfiguration) Host() string {
		return this.Get(CnfKvRedisHost, CnfKvRedisHostDefault)
}

func (this *RedisConfiguration) Collection() string {
		return this.Get(CnfKvRedisCollection, CnfKvRedisCollectionDefault)
}

func (this *RedisConfiguration) Port() int {
		return this.GetInt(CnfKvRedisPort, CnfKvRedisPortDefault)
}

func (this *RedisConfiguration) Conn() string {
		return fmt.Sprintf("%s:%d", this.Host(), this.Port())
}

func (this *RedisConfiguration) DbNum() int {
		return this.GetInt(CnfKvRedisDbNum, CnfKvRedisDbNumDefault)
}

func (this *RedisConfiguration) GetInt(key string, def ...int) int {
		if len(def) == 0 {
				def = append(def, 0)
		}
		kv := mapper.Mapper(this.KvCnf)
		return kv.GetInt(this.Key(key), def...)
}

func (this *RedisConfiguration) Key(key string) string {
		if key == "" {
				return ""
		}
		if this.ScopeName == "" {
				this.ScopeName = CnfKvRedisGlobal
		}
		return this.ScopeName + "." + key
}

func (this *RedisConfiguration) Destroy() {
		this.CacheConfiguration.Destroy()
}
