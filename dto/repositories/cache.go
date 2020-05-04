package repositories

import (
		"github.com/WebGameLinux/cms/dto/configuration"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/astaxie/beego/cache"
		_ "github.com/astaxie/beego/cache/memcache"
		_ "github.com/astaxie/beego/cache/redis"
		_ "github.com/astaxie/beego/cache/ssdb"
		"time"
)

type CacheRepositoryManager struct {
		repos     map[string]cache.Cache
		configure mapper.Mapper
}

type CacheManager interface {
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

func NewCacheManager(config ...mapper.Mapper) CacheManager {
		if cacheRepository == nil {
				cacheRepository = createManager(config...)
				// 注册 缓存
				cacheRepository.Register(DriverStoreDefault, driver(DriverStoreDefault))
				cacheRepository.Register(DriverStoreRedis, driver(DriverStoreRedis))
		}
		return cacheRepository
}

func driver(name string) cache.Cache {
		c, _ := cache.NewCache(name, configuration.GetCacheConfig(name))
		return c
}

func GetCacheManager(config ...mapper.Mapper) *CacheRepositoryManager {
		m := NewCacheManager(config...)
		manager, _ := m.(*CacheRepositoryManager)
		return manager
}

func createManager(config ...mapper.Mapper) *CacheRepositoryManager {
		var manager = new(CacheRepositoryManager)
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

func (this *CacheRepositoryManager) Get(key string) interface{} {
		return this.Store().Get(key)
}

func (this *CacheRepositoryManager) GetMulti(keys []string) []interface{} {
		return this.Store().GetMulti(keys)
}

func (this *CacheRepositoryManager) Put(key string, val interface{}, timeout time.Duration) error {
		return this.Store().Put(key, val, timeout)
}

func (this *CacheRepositoryManager) Delete(key string) error {
		return this.Store().Delete(key)
}

func (this *CacheRepositoryManager) Incr(key string) error {
		return this.Store().Incr(key)
}

func (this *CacheRepositoryManager) Decr(key string) error {
		return this.Store().Decr(key)
}

func (this *CacheRepositoryManager) IsExist(key string) bool {
		return this.Store().IsExist(key)
}

func (this *CacheRepositoryManager) ClearAll() error {
		return this.Store().ClearAll()
}

func (this *CacheRepositoryManager) StartAndGC(config string) error {
		return this.Store().StartAndGC(config)
}

func (this *CacheRepositoryManager) Store(name ...string) cache.Cache {
		if len(name) == 0 {
				name = append(name, StoreDefault)
		}
		var key = name[0]
		if store, ok := this.repos[key]; ok && store != nil {
				return store
		}
		return this.getDefault()
}

func (this *CacheRepositoryManager) getDefault() cache.Cache {
		if store, ok := this.repos[StoreDefault]; ok && store != nil {
				return store
		}
		c, err := cache.NewCache(DriverStoreDefault, configuration.GetCacheConfig(DriverStoreDefault))
		if err != nil {
				this.Register(StoreDefault, c)
		}
		return c
}

func (this *CacheRepositoryManager) Register(name string, instance cache.Cache) {
		if name == "" || instance == nil {
				return
		}
		this.repos[name] = instance
}
