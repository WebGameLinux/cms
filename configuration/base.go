package configuration

import "fmt"

type ConfigureDto struct {
		KvCnf     map[string]interface{}
}

type BaseConfiguration interface {
		Get(key string, def ...string) string
		GetInt(key string, def ...int) int
		fmt.Stringer
		Destroy()
}

