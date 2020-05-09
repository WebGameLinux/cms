package redis

import (
		"github.com/go-redis/redis"
		"log"
		"time"
)

type RedisOptions struct {
		Network string `json:"network"`
		// host:port address.
		Addr string `json:"addr"`
		// Optional password. Must match the password specified in the
		// requirepass server configuration option.
		Password string `json:"password"`
		// Database to be selected after connecting to the server.
		DB int `json:"db"`
		// Maximum number of retries before giving up.
		// Default is to not retry failed commands.
		MaxRetries int `json:"max_retries"`
		// Minimum backoff between each retry.
		// Default is 8 milliseconds; -1 disables backoff.
		MinRetryBackoff time.Duration `json:"min_retry_backoff"`
		// Maximum backoff between each retry.
		// Default is 512 milliseconds; -1 disables backoff.
		MaxRetryBackoff time.Duration `json:"max_retry_backoff"`
		// Dial timeout for establishing new connections.
		// Default is 5 seconds.
		DialTimeout time.Duration `json:"dial_timeout"`
		// Timeout for socket reads. If reached, commands will fail
		// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
		// Default is 3 seconds.
		ReadTimeout time.Duration `json:"read_timeout"`
		// Timeout for socket writes. If reached, commands will fail
		// with a timeout instead of blocking.
		// Default is ReadTimeout.
		WriteTimeout time.Duration `json:"write_timeout"`

		// Maximum number of socket connections.
		// Default is 10 connections per every CPU as reported by runtime.NumCPU.
		PoolSize int `json:"pool_size"`
		// Minimum number of idle connections which is useful when establishing
		// new connection is slow.
		MinIdleConns int `json:"min_idle_conns"`
		// Connection age at which client retires (closes) the connection.
		// Default is to not close aged connections.
		MaxConnAge time.Duration `json:"max_conn_age"`
		// Amount of time client waits for connection if all connections
		// are busy before returning an error.
		// Default is ReadTimeout + 1 second.
		PoolTimeout time.Duration `json:"pool_timeout"`
		// Amount of time after which client closes idle connections.
		// Should be less than server's timeout.
		// Default is 5 minutes. -1 disables idle timeout check.
		IdleTimeout time.Duration `json:"idle_timeout"`
		// Frequency of idle checks made by idle connections reaper.
		// Default is 1 minute. -1 disables idle connections reaper,
		// but idle connections are still discarded by the client
		// if IdleTimeout is set.
		IdleCheckFrequency time.Duration `json:"idle_check_frequency"`
}

func NewRedisOptions(data string) *RedisOptions {
		var option = new(RedisOptions)
		if err := option.UnmarshalJSON([]byte(data)); err != nil {
				log.Println(err)
		}
		return option
}

func (this *RedisOptions) Options() *redis.Options {
		var opt = new(redis.Options)
		opt.Password = this.Password
		opt.Network = this.Network
		opt.Addr = this.Addr
		opt.DB = this.DB
		opt.DialTimeout = this.DialTimeout
		opt.MaxConnAge = this.MaxConnAge
		opt.PoolSize = this.PoolSize
		opt.WriteTimeout = this.WriteTimeout
		opt.ReadTimeout = this.ReadTimeout
		opt.MinIdleConns = this.MinIdleConns
		opt.MaxRetries = this.MaxRetries
		opt.MaxRetryBackoff = this.MaxRetryBackoff
		opt.MinRetryBackoff = this.MinRetryBackoff
		opt.PoolTimeout = this.PoolTimeout
		opt.IdleTimeout = this.IdleTimeout
		opt.IdleCheckFrequency = this.IdleCheckFrequency
		return opt
}
