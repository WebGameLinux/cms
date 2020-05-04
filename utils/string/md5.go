package string

import (
		"crypto/md5"
		"fmt"
)

func Md5(string string) string {
		h := md5.New()
		h.Write([]byte(string))
		return fmt.Sprintf("%x", h.Sum(nil))
}

// Md5(reflect.TypeOf(options).PkgPath())