package types

import (
		"fmt"
		"github.com/WebGameLinux/cms/utils/mapper"
		"sync"
)

type ProviderManagerDto struct {
		Container sync.Map
}

type BaseProviderManager struct {
		ProviderManagerDto
		Booted                         bool
		BootHandler                    func(self *BaseProviderManager)
		ResolverHandler                func(key string, self *BaseProviderManager) (Provider, bool)         // 自定义解析器
		RegisterHandler                func(key string, instance Provider, self *BaseProviderManager)       // 自定义注册器
		RegisterProviderFactoryHandler func(key string, factory func() Provider, self *BaseProviderManager) // 自定义工厂注册器
}

func NewBaseProviderManager() ProviderManagerContainer {
		return new(BaseProviderManager)
}

func ProviderManager2BaseProviderManager(p ProviderManagerContainer) *BaseProviderManager {
		if base, ok := p.(*BaseProviderManager); ok {
				return base
		}
		return nil
}

func (this *BaseProviderManager) Resolver(key string) (Provider, bool) {
		if key == "" {
				return nil, false
		}
		if !this.Booted {
				this.Boot()
		}
		if this.ResolverHandler != nil {
				return this.ResolverHandler(key, this)
		}
		v, ok := this.Container.Load(key)
		if !ok {
				return nil, false
		}
		if p, ok := v.(Provider); ok {
				return p, ok
		}
		if p, ok := v.(func() Provider); ok {
				return p(), ok
		}
		return nil, false
}

func (this *BaseProviderManager) Boot() {
		if this.Booted {
				return
		}
		if this.BootHandler != nil {
				this.BootHandler(this)
		}
		this.Booted = true
}

func (this *BaseProviderManager) Register(key string, instance Provider) {
		if key == "" || instance == nil {
				return
		}
		if !this.Booted {
				this.Boot()
		}
		if this.RegisterHandler != nil {
				this.RegisterHandler(key, instance, this)
				return
		}
		this.Container.Store(key, instance)
}

func (this *BaseProviderManager) RegisterProviderFactory(key string, factory func() Provider) {
		if key == "" || factory == nil {
				return
		}
		if !this.Booted {
				this.Boot()
		}
		if this.RegisterProviderFactoryHandler != nil {
				this.RegisterProviderFactoryHandler(key, factory, this)
				return
		}
		this.Container.Store(key, factory)
}

func (this *BaseProviderManager) Lists() mapper.Mapper {
		var m = make(mapper.Mapper)
		this.Container.Range(func(key, value interface{}) bool {
				if k, ok := key.(string); ok && key != "" && value != nil {
						m[k] = value
				}
				return true
		})
		return m
}

func (this *BaseProviderManager) Count() int {
		var length = 0
		this.Container.Range(func(key, value interface{}) bool {
				length++
				return true
		})
		return length
}

func (this *BaseProviderManager) Dump() {
		fmt.Println("{")
		this.Container.Range(func(key, value interface{}) bool {
				if k, ok := key.(string); ok && key != "" && value != nil {
						fmt.Printf(`"%s":<Provider>%+v`+"\n", k, value)
				}
				return true
		})
		fmt.Println("}")
}

func (this *BaseProviderManager) Destroy() {
		this.Container.Range(func(key, value interface{}) bool {
				this.Container.Delete(key)
				return true
		})
}

func (this *BaseProviderManager) Delete(key string) {
		this.Container.Delete(key)
}
