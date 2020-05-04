package types

import (
		"github.com/WebGameLinux/cms/utils/mapper"
		"sync"
)

// 相关服务提供器
type Provider interface {
		Get(string) (interface{}, bool)            // 获取相关 Provider method
		Set(string, interface{})                   // 设置相关 Provider methods
		Invoke(...interface{}) Provider            // 新构造一个 Provider
		Bind(string, interface{}) Provider         // 绑定相关参数
		Arg(string) (interface{}, bool)            // 获取相关参数
		Name() string                              // 服务名
		Register(manager ProviderManagerContainer) // 注册到管理器
}

type ProviderArgs struct {
		Name      string
		Container *sync.Map
		Args      mapper.Mapper
}

// 服务提供器管理
type ProviderManagerContainer interface {
		Boot()
		Resolver(string) (Provider, bool)
		Register(string, Provider)
		RegisterProviderFactory(string, func() Provider)
}

type BaseProviderDto struct {
		Container sync.Map      // 方法
		Name      string        // 服务名
		Args      mapper.Mapper // 参数
}

type BaseProvider struct {
		BaseProviderDto
}

func NewBaseProvider(args ...interface{}) Provider {
		if len(args) == 0 {
				return nil
		}
		var p *BaseProvider
		for _, v := range args {
				if n, ok := v.(string); ok && n != "" {
						if p == nil {
								p = new(BaseProvider)
						}
						if p.BaseProviderDto.Name == "" {
								p.BaseProviderDto.Name = n
						}
						continue
				}
				if arg, ok := v.(*ProviderArgs); ok {
						if p == nil {
								if arg.Name == "" {
										continue
								}
								p = new(BaseProvider)
						}
						// 容器
						if arg.Container != nil {
								p.Container = *arg.Container
						}
						// 设置名字
						if arg.Name != "" && p.BaseProviderDto.Name == "" {
								p.BaseProviderDto.Name = arg.Name
						}
						// 绑定参数
						if arg.Args != nil && len(arg.Args) > 0 {
								for k, v := range arg.Args {
										p.Bind(k, v)
								}
						}
				}
		}
		return p
}

func Provider2BaseProvider(p Provider) *BaseProvider {
		if base, ok := p.(*BaseProvider); ok {
				return base
		}
		return nil
}

func (this *BaseProvider) Get(method string) (interface{}, bool) {
		return this.Container.Load(method)
}

func (this *BaseProvider) Name() string {
		return this.BaseProviderDto.Name
}

func (this *BaseProvider) Set(method string, handler interface{}) {
		this.Container.Store(method, handler)
}

func (this *BaseProvider) Invoke(args ...interface{}) Provider {
		return NewBaseProvider(args...)
}

func (this *BaseProvider) Bind(key string, v interface{}) Provider {
		this.Args[key] = v
		return this
}

func (this *BaseProvider) Arg(key string) (interface{}, bool) {
		if v, ok := this.Args[key]; ok {
				return v, ok
		}
		return nil, false
}

func (this *BaseProvider) Register(manager ProviderManagerContainer) {
		manager.Register(this.Name(), this)
}

func (this *BaseProvider) Destroy() {
		this.Container.Range(func(key, value interface{}) bool {
				this.Container.Delete(key)
				return true
		})
}

func (this *BaseProvider) Delete(key string) {
		this.Container.Delete(key)
}
