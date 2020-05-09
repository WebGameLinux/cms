package repositories

import "sync"

var (
		lock            sync.Once
)

func init() {
		lock.Do(func() {
				GetUserRepository()
		})
}
