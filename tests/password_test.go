package test

import (
		"fmt"
		. "github.com/WebGameLinux/cms/utils/string"
		. "github.com/smartystreets/goconvey/convey"
		"testing"
)

func TestPassword(t *testing.T) {
		var text = "hello-word-golang"
		// 32.1.c77ecd55df3ecd97bb9aded010806a6d97d152edd4901ab719170708aa3e87d1.32768.8.MTI4
		// 32.1.c267402682ba56024de2841d44c73370ff8092ac00534571c07e77ae8c564688.32768.8.MTIzNw==
		Convey("密码函数测试", t, func() {
				pass := PasswordHash(text)
				So(pass, ShouldNotBeNil)
				So(pass, ShouldNotBeEmpty)
				Convey("验证密码", func() {
						So(PasswordVerify(text, pass), ShouldEqual, true)
				})
		})
}

func MakeN(size int, macs ...int) {
		var (
				i   = 0
				arr []int
				min = 100
				max = 40000
		)
		if len(macs) > 0 && macs[0] > min {
				max = macs[0]
		}
		for i = min; i < max; i++ {
				if i&(i-1) == 0 && i > 1 {
						arr = append(arr, i)
				}
				if len(arr) >= size {
						break
				}
		}
		fmt.Println(arr)
}

func TestMakeN(t *testing.T) {
		MakeN(10,100000)
}

func BenchmarkPassword(b *testing.B) {
		var text = "hello-word-golang"
		for i := 0; i < b.N; i++ {
				Convey("密码函数测试", b, func() {
						pass := PasswordHash(text)
						So(pass, ShouldNotBeNil)
						So(pass, ShouldNotBeEmpty)
						//	fmt.Println(pass)
						Convey("验证密码", func() {
								So(PasswordVerify(text, pass), ShouldEqual, true)
								So(PasswordVerify(text+"123", pass), ShouldEqual, false)
						})
				})
		}
}
