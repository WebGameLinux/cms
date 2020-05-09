package test

import (
		"fmt"
		"github.com/WebGameLinux/cms/dto/services"
		"github.com/WebGameLinux/cms/libs/redis"
		. "github.com/smartystreets/goconvey/convey"
		"testing"
		"time"
)

func BenchmarkNewRedisService(b *testing.B) {

}

func TestNewRedisService(t *testing.T) {
		var service = services.NewRedisService()
		service.Get("test")
}

func TestRedisGetInstance(t *testing.T) {
		var (
				key      = "user:token:1"
				instance = redis.GetInstance(redis.TokenDb)
				value    = fmt.Sprintf("%d", time.Now().Unix())
		)

		Convey("redis set test", t, func() {
				v, err := instance.Set(key, value, 0).Result()
				So(err, ShouldBeNil)
				So(v, ShouldEqual, "OK")
				v, err = instance.Get(key).Result()
				So(err, ShouldBeNil)
				So(v, ShouldEqual, value)
				t, err := instance.TTL(key).Result()
				So(err, ShouldBeNil)
				So(t, ShouldEqual, -1*time.Second)

				instance.Expire(key, 5*time.Minute)
				t, err = instance.TTL(key).Result()
				So(err, ShouldBeNil)
				So(t, ShouldNotEqual, time.Duration(0))
		})
}
