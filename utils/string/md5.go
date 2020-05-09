package string

import (
		"crypto/md5"
		"fmt"
		"github.com/WebGameLinux/cms/utils/beego"
		"time"
)

func Md5(string string) string {
		h := md5.New()
		h.Write([]byte(string))
		return fmt.Sprintf("%x", h.Sum(nil))
}

func SeqId(table string) string {
		return Md5(beego.GetKvString("node", "default") + table + time.Now().String())
}

// Md5(reflect.TypeOf(options).PkgPath())
