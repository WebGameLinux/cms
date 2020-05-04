package test

import (
		"github.com/WebGameLinux/cms/dto/services"
		. "github.com/smartystreets/goconvey/convey"
		"testing"
)

func TestNewSuccessResult(t *testing.T) {
		Convey("创建结果结构对象", t, func() {
				res := services.NewSuccessResult(services.Map{"name": "user", "array": []int{1, 2, 4, 5}})
				So(res, ShouldNotBeNil)
				Convey("result->map转换", func() {
						m := res.Mapper()
						So(m, ShouldNotBeNil)
						So(m["data"], ShouldNotBeNil)
						So(m["code"], ShouldEqual, 0)
						So(m["message"], ShouldEqual, "ok")
						So(m["error"], ShouldBeNil)
				})
		})
}
