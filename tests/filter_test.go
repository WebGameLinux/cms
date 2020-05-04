package test

import (
		. "github.com/WebGameLinux/cms/utils/array"
		. "github.com/smartystreets/goconvey/convey"
		"testing"
)

func TestStrArrayFilter(t *testing.T) {
		var (
				regexpLeft = 5
				arr        = []string{"", "1231", "1", "sf", "\n", "\t\n", "\t\t", " ", " ", "user",}
		)
		Convey("正则过滤空字符串(默认)", t, func() {
				arrNew := StrArrayRegexpFilter(arr)
				So(len(arrNew), ShouldEqual, regexpLeft)
		})
}

func TestContains(t *testing.T) {
		var (
				arr = []string{"", "1231", "1", "sf", "\n", "\t\n", "\t\t", " ", " ", "user",}
				sub = " "
		)
		Convey("查找字符串数组中是否含子字符", t, func() {
				So(Contains(arr, sub), ShouldEqual, true)
		})
}

func TestCount(t *testing.T) {
		var (
				arr    = []string{"", "1231", "1", "sf", "\n", "\t\n", "\t\t", " ", " ", "user",}
				sub    = " "
				subNum = 2
		)
		Convey("查找字符串数组中是否含子字符", t, func() {
				So(Count(arr, sub), ShouldEqual, subNum)
				So(Count(arr, ""), ShouldEqual, 1)
		})
}
