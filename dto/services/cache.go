package services

import (
		"github.com/WebGameLinux/cms/configuration"
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/astaxie/beego/cache"
		_ "github.com/astaxie/beego/cache/memcache"
		_ "github.com/astaxie/beego/cache/redis"
		_ "github.com/astaxie/beego/cache/ssdb"
		"time"
)

type CacheManagerService struct {
		repos     map[string]cache.Cache
		configure mapper.Mapper
}

type CacheManagerServiceInterface interface {
		cache.Cache
		Store(name ...string) cache.Cache
		Register(name string, instance cache.Cache)
}

const (
		StoreDefault        = "default"
		DriverStoreFile     = "file"
		DriverStoreRedis    = "redis"
		DriverStoreMemcache = "memcache"
		DriverStoreMemory   = "memory"
		DriverStoreSsd      = "ssd"
		DriverStoreDefault  = DriverStoreFile
)

var cacheManager CacheManagerServiceInterface

func NewCacheManager(config ...mapper.Mapper) CacheManagerServiceInterface {
		if cacheManager == nil {
				cacheManager = createManager(config...)
				// 注册 缓存
				cacheManager.Register(DriverStoreDefault, driver(DriverStoreDefault))
				cacheManager.Register(DriverStoreRedis, driver(DriverStoreRedis))
				cacheManager.Register(enums.TokenStore, CreateRedisService(enums.TokenStore))
		}
		return cacheManager
}

func driver(name string) cache.Cache {
		c, _ := cache.NewCache(name, configuration.GetCacheConfig(name))
		return c
}

func GetCacheManagerService(config ...mapper.Mapper) *CacheManagerService {
		m := NewCacheManager(config...)
		manager, _ := m.(*CacheManagerService)
		return manager
}

func createManager(config ...mapper.Mapper) *CacheManagerService {
		var manager = new(CacheManagerService)
		if len(config) == 0 {
				config = append(config, GetCacheConfigure())
		}
		manager.configure = config[0]
		manager.repos = make(map[string]cache.Cache)
		return manager
}

// 获取缓存配置
func GetCacheConfigure() mapper.Mapper {
		var m = make(mapper.Mapper)
		return m
}

func (this *CacheManagerService) Get(key string) interface{} {
		return this.Store().Get(key)
}

func (this *CacheManagerService) GetMulti(keys []string) []interface{} {
		return this.Store().GetMulti(keys)
}

func (this *CacheManagerService) Put(key string, val interface{}, timeout time.Duration) error {
		return this.Store().Put(key, val, timeout)
}

func (this *CacheManagerService) Delete(key string) error {
		return this.Store().Delete(key)
}

func (this *CacheManagerService) Incr(key string) error {
		return this.Store().Incr(key)
}

func (this *CacheManagerService) Decr(key string) error {
		return this.Store().Decr(key)
}

func (this *CacheManagerService) IsExist(key string) bool {
		return this.Store().IsExist(key)
}

func (this *CacheManagerService) ClearAll() error {
		return this.Store().ClearAll()
}

func (this *CacheManagerService) StartAndGC(config string) error {
		return this.Store().StartAndGC(config)
}

func (this *CacheManagerService) Store(name ...string) cache.Cache {
		if len(name) == 0 {
				name = append(name, StoreDefault)
		}
		var key = name[0]
		if store, ok := this.repos[key]; ok && store != nil {
				return store
		}
		return this.getDefault()
}

func (this *CacheManagerService) getDefault() cache.Cache {
		if store, ok := this.repos[StoreDefault]; ok && store != nil {
				return store
		}
		c, err := cache.NewCache(DriverStoreDefault, configuration.GetCacheConfig(DriverStoreDefault))
		if err != nil {
				this.Register(StoreDefault, c)
		}
		return c
}

func (this *CacheManagerService) Register(name string, instance cache.Cache) {
		if name == "" || instance == nil {
				return
		}
		this.repos[name] = instance
}
