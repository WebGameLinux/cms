package configuration

type CacheConfiguration struct {
		ConfigureDto
		ScopeName string
}

type CacheConfigure interface {
		BaseConfiguration
		Scope(name string) CacheConfigure
		Args() []interface{}
}

func (this *CacheConfiguration) Destroy() {
		this.ConfigureDto.KvCnf = nil
}

func GetCacheConfig(name string,scope ...string) string {
		switch name {
		case "file":
				return FileCnfKvString()
		case "redis":
				return RedisCnfString()
		default:
				return FileCnfKvString()
		}
}
