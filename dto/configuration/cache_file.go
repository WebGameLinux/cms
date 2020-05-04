package configuration

import (
		"fmt"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/astaxie/beego/config"
)

type FileCacheConfiguration struct {
		CacheConfiguration
}

var fileKvCnf *FileCacheConfiguration

const (
		CnfKvFileCacheEmbedExpiry           = "EmbedExpiry"
		CnfKvFileCacheDirectoryLevel        = "DirectoryLevel"
		CnfKvFileCacheCachePath             = "CachePath"
		CnfKvFileCacheFileSuffix            = "FileSuffix"
		CnfKvFileCacheCachePathDefault      = "runtime/cache"
		CnfKvFileCacheGlobal                = "cache.file"
		CnfKvFileCacheFileSuffixDefault     = "app_cache"
		CnfKvFileCacheDirectoryLevelDefault = 3
		CnfKvFileCacheEmbedExpiryDefault    = true
		StringFileCacheTemplate             = `{"CachePath":"%s","FileSuffix":"%s","DirectoryLevel":"%d","EmbedExpiry":"%s"}`
)

type FileCacheConfigure interface {
		CacheConfigure
		GetBool(string, ...bool) bool
		CachePath() string
		FileSuffix() string
		DirectoryLevel() int
		EmbedExpiry() bool
}

func FileCnfKvString() string {
		return fileKvCnf.String()
}

func FileCnfScope(name string) CacheConfigure {
		return fileKvCnf.Scope(name)
}

func GetFileCacheKvCnf(Properties ...interface{}) *FileCacheConfiguration {
		if fileKvCnf == nil {
				fileKvCnf = new(FileCacheConfiguration)
				fileKvCnf.init()
		}
		if len(Properties) > 0 {
				kv := Properties[0]
				switch kv.(type) {
				case mapper.Mapper:
						if m, ok := kv.(mapper.Mapper); ok {
								fileKvCnf.KvCnf = m
						}
				case config.Configer:
						if m, ok := kv.(config.Configer); ok {
								fileKvCnf.KvCnf = AppConfig2Map(m, CnfKvFileCacheGlobal)
						}
				default:
						fileKvCnf.KvCnf = map[string]interface{}{}
				}
		}
		return fileKvCnf
}

func (this *FileCacheConfiguration) init() {
		if this.KvCnf == nil {
				this.KvCnf = make(map[string]interface{})
		}
		if this.ScopeName == "" {
				this.ScopeName = CnfKvFileCacheGlobal
		}
}

func (this *FileCacheConfiguration) String() string {
		return fmt.Sprintf(StringFileCacheTemplate, this.Args()...)
}

func (this *FileCacheConfiguration) Args() []interface{} {
		return []interface{}{
				this.CachePath(),
				this.FileSuffix(),
				this.DirectoryLevel(),
				this.EmbedExpiry(),
		}
}

func (this *FileCacheConfiguration) Scope(name string) CacheConfigure {
		kv := new(FileCacheConfiguration)
		kv.init()
		kv.ScopeName = this.ScopeName + "." + name
		kv.KvCnf = this.KvCnf
		return kv
}

func (this *FileCacheConfiguration) Get(key string, def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		kv := mapper.Mapper(this.KvCnf)
		return kv.Get(this.key(key), def...)
}

func (this *FileCacheConfiguration) CachePath() string {
		return this.Get(CnfKvFileCacheCachePath, CnfKvFileCacheCachePathDefault)
}

func (this *FileCacheConfiguration) FileSuffix() string {
		return this.Get(CnfKvFileCacheFileSuffix, CnfKvFileCacheFileSuffixDefault)
}

func (this *FileCacheConfiguration) EmbedExpiry() bool {
		return this.GetBool(CnfKvFileCacheEmbedExpiry, CnfKvFileCacheEmbedExpiryDefault)
}

func (this *FileCacheConfiguration) DirectoryLevel() int {
		return this.GetInt(CnfKvFileCacheDirectoryLevel, CnfKvFileCacheDirectoryLevelDefault)
}

func (this *FileCacheConfiguration) GetInt(key string, def ...int) int {
		if len(def) == 0 {
				def = append(def, 0)
		}
		kv := mapper.Mapper(this.KvCnf)
		return kv.GetInt(this.key(key), def...)
}

func (this *FileCacheConfiguration) GetBool(key string, def ...bool) bool {
		if len(def) == 0 {
				def = append(def, false)
		}
		kv := mapper.Mapper(this.KvCnf)
		return kv.GetBool(this.key(key), def...)
}

func (this *FileCacheConfiguration) key(key string) string {
		if key == "" {
				return ""
		}
		if this.ScopeName == "" {
				this.ScopeName = CnfKvFileCacheGlobal
		}
		return this.ScopeName + "." + key
}

func (this *FileCacheConfiguration) Destroy() {
		this.CacheConfiguration.Destroy()
}
