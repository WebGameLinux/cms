package test

import (
		"fmt"
		"github.com/WebGameLinux/cms/utils/configure"
		"testing"
		"time"
)

var file = `file://conf?mode=development&appendSystemEnv=false`
var zookeeper = `zookeeper://?connStr=182.61.3.40:16200`

func TestNewConfiguration(t *testing.T) {
		var url = zookeeper
		var config = configure.GetConfigure(url)
		// config.Set("app_url","www.api.com")
	 //  config.Set(".emails","122@qq.com,12323@163.com")
		fmt.Println(config.Get(".app"))
		fmt.Println(config.Get(".user"))
		fmt.Println(config.Strings(".user.emails"))
		kvs:=config.Scope("api")
		fmt.Println(kvs)
		time.Sleep(50 * time.Second)
}
