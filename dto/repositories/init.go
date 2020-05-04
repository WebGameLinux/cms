package repositories

import "sync"

var (
		lock            sync.Once
		cacheRepository CacheManager
)

func init() {
		lock.Do(func() {
				NewCacheManager()
				GetUserRepository()
		})
}
