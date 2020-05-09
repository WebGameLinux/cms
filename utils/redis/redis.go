package redis

import (
		"github.com/go-redis/redis"
		"sync"
)

type Container struct {
		clients map[string]*redis.Client
		sync.RWMutex
}

func NewContainer() *Container {
		container := new(Container)
		container.clients = make(map[string]*redis.Client)
		return container
}

var container *Container

func GetRedis(name ...string) *redis.Client {
		if len(name) == 0 {
				name = append(name, "default")
		}
		return GetRedisContainer().Get(name[0])
}

func GetRedisContainer() *Container {
		if container == nil {
				container = NewContainer()
		}
		return container
}

func (this *Container) Get(name string) *redis.Client {
		return this.clients[name]
}

func (this *Container) Add(name string, config string) {
		this.Lock()
		defer this.Unlock()
		if this.has(name) {
				return
		}
		var opts = this.resolver(config)
		this.clients[name] = redis.NewClient(opts)
}

func (this *Container) has(name string) bool {
		if _, ok := this.clients[name]; ok {
				return true
		}
		return false
}

func (this *Container) resolver(config string) *redis.Options {
		var opt = NewRedisOptions(config)
		return opt.Options()
}
