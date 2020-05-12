package test

import (
		"fmt"
		"github.com/WebGameLinux/cms/libs/filesystem/local"
		"os"
		"testing"
)

func TestGetFile(t *testing.T) {
		//	fmt.Println(local.EmptySize.ParseInt(1))
		//	fmt.Println(local.EmptySize.Parse("9990.2886GB"))
		//fmt.Println(local.GetFileSizeFormat("./", "%s%s"))
	//	fmt.Println(local.EmptySizeNum.ParseFloatN(local.GB + 1))
	//	fmt.Println(local.EmptySizeNum.Parse2Num("1.00GB").Prev(4))
	//	fmt.Println(local.GetTwoUtilMultiple("GB","ZB"))
		var fs *os.File
		disk :=local.GetDisk()
		disk.Proxy("chRoot","./tests")
		fd:= disk.Open("config_test.go")
		if fd !=nil {
				fs = fd.(*os.File)
				fmt.Println(fs.Name())
		}
	//	fmt.Println(disk.GetError())
}


