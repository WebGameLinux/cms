package string

import (
		"crypto/md5"
		"encoding/hex"
		"fmt"
		"github.com/WebGameLinux/cms/utils/beego"
		"io"
		"os"
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
// 文件hash
func FileHash(filename string) string {
		var (
				err  error
				file *os.File
				info os.FileInfo
		)
		if info, err = os.Stat(filename); err != nil {
				return ""
		}
		if info.IsDir() {
				return ""
		}
		file, err = os.OpenFile(filename, os.O_RDWR, os.ModePerm)
		if err != nil {
				return ""
		}
		hash := md5.New()
		_, err = io.Copy(hash, file)
		if err != nil {
				return ""
		}
		return hex.EncodeToString(hash.Sum(nil))
}
