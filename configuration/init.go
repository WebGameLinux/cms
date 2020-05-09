package configuration

import (
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/WebGameLinux/cms/utils/reflects"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/config"
		"sync"
)

var (
		lock sync.Once
)

func init() {
		lock.Do(func() {
				GetRedisKvCnf(AppConfig2Map(beego.AppConfig, CnfKvRedisGlobal))
				GetFileCacheKvCnf(AppConfig2Map(beego.AppConfig, CnfKvFileCacheGlobal))
		})
}

func AppConfig2Map(config config.Configer, section string) mapper.Mapper {
		m, _ := config.GetSection(section)
		kv, _ := reflects.Map2Mapper(m)
		return kv
}
