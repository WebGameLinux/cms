package redis

const (
		QueueDb     = "queue"      // 队列库
		LockDb      = "lock"       // 锁存储库
		TokenDb     = "token"      // 临牌缓存
		TempDb      = "default"    // 临时存储
		CacheDb     = "cache"      // 数据库缓存
		StatisticDb = "statistics" // 分析统计缓存

)

var NameTypeMap map[string]int

func initNames() {
		NameTypeMap = map[string]int{
				QueueDb:     0,
				LockDb:      1,
				TokenDb:     2,
				TempDb:      3,
				StatisticDb: 4,
				CacheDb:     5,
		}
}

func getNameMap() map[string]int {
		if len(NameTypeMap) == 0 || NameTypeMap == nil {
				initNames()
		}
		return NameTypeMap
}

func getDbByName(name string) int {
		m := getNameMap()
		if db, ok := m[name]; ok {
				return db
		}
		return m[TempDb]
}

func getDbName(db int) string {
		m := getNameMap()
		for name, n := range m {
				if n == db {
						return name
				}
		}
		return TempDb
}
