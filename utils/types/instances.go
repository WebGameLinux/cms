package types

import "sync"

var (
		once                   sync.Once
		providerMangerInstance ProviderManagerContainer
)

func init() {
		once.Do(func() {
				if providerMangerInstance == nil {
						providerMangerInstance = NewBaseProviderManager()
				}
		})
}

func GetProviderMangerInstance() ProviderManagerContainer {
		return providerMangerInstance
}

func Resolver(key string) (Provider, bool) {
		return providerMangerInstance.Resolver(key)
}

func Register(key string, instance Provider) {
		providerMangerInstance.Register(key, instance)
}

func RegisterProviderFactory(key string, factory func() Provider) {
		providerMangerInstance.RegisterProviderFactory(key, factory)
}
