package beego

import (
		"errors"
		"github.com/WebGameLinux/cms/utils/array"
		"github.com/astaxie/beego"
		"regexp"
		"strings"
)

var (
		cacheExpress  = make(map[string]interface{})
		expressVarReg = regexp.MustCompile(`(\$\{([a-zA-Z|_.]+)\})`)
)

// 运行模式
func RunMode(refresh ...bool) string {
		mode := beego.BConfig.RunMode
		if mode == "" {
				return "dev"
		}
		return strings.Trim(mode, " ")
}

// 获取环境配置
func GetKv(ty, key string, def interface{}) (interface{}, error) {
		mode := RunMode()
		if key == "" {
				return def, errors.New("empty key")
		}
		if strings.Contains(key, mode+"::") {
				key = strings.Replace(key, mode+"::", "", -1)
		}
		return beego.GetConfig(UcFirst(ty), key, def)
}

// 首字母大写
func UcFirst(str string) string {
		arr := []rune(str)
		first := string(arr[0])
		upperFirst := strings.ToUpper(first)
		if upperFirst != first {
				arr = append([]rune(upperFirst), arr[1:]...)
				return string(arr)
		}
		return str
}

// 是否对应环境
func IsThatRunMode(mode string) bool {
		return strings.ToLower(RunMode()) == strings.ToLower(mode)
}

// 获取字符串带 runmode
func GetKvString(key string, def ...string) string {
		if len(def) <= 0 {
				def = append(def, "")
		}
		v, _ := GetKv("string", key, def[0])
		return v.(string)
}

// 获取整形配置
func GetKvInt(key string, def ...int) int {
		if len(def) <= 0 {
				def = append(def, 0)
		}
		v, _ := GetKv("Int", key, def[0])
		return v.(int)
}

// 获取带前缀的表名
func GetTable(name string) string {
		return name
}

// 解析变量
func GetParsedKvStr(key string, def ...string) string {
		v := GetKvString(key, def...)
		if v == "" {
				return v
		}
		return ParseVarExpressVar(v)
}

// 获取配置bool类型
func GetKvBool(key string, def ...bool) bool {
		if len(def) == 0 {
				def = append(def, false)
		}
		if v, err := GetKv("Bool", key, def[0]); err == nil {
				return v.(bool)
		}
		return def[0]
}

// 解析变量表达式
func ParseVarExpressVar(v string) string {
		var (
				i   = 0
				num int
				val string
		)
		vars := expressVarReg.FindAllStringSubmatch(v, -1)
		num = len(vars)

		if num <= 0 {
				return v
		}
		for ; i < num; i++ {
				if len(vars[0]) < 3 {
						continue
				}
				key := vars[i][1]
				// 缓存中获取
				if value, ok := GetVarInCacheMap(key); ok {
						v = strings.ReplaceAll(v, key, value.(string))
						continue
				}
				varArg := strings.SplitN(vars[i][2], "||", -1)
				varArgc := len(varArg)
				if varArgc == 0 {
						continue
				}

				if varArgc > 0 && varArgc < 2 {
						val = GetKvString(varArg[0])
				}
				if varArgc >= 2 {
						val = GetKvString(varArg[0], varArg[1])
				}
				if val == "" {
						continue
				}
				// 多重表达式
				if HasVarExpress(val) {
						val = ParseVarExpressVar(val)
				}
				setVarInCacheMap(key, val)
				v = strings.ReplaceAll(v, key, val)
		}
		return v
}

func GetVarInCacheMap(name string) (interface{}, bool) {
		if v, ok := cacheExpress[name]; ok {
				return v, true
		}
		return nil, false
}

// 设置缓存
func setVarInCacheMap(name string, val interface{}) {
		cacheExpress[name] = val
}

// 清理缓存
func DelVarInCacheMap(name ...string) int {
		var i = 0
		if len(name) == 0 {
				for k, _ := range cacheExpress {
						delete(cacheExpress, k)
						i++
				}
		} else {
				for _, key := range name {
						if _, ok := cacheExpress[key]; ok {
								delete(cacheExpress, key)
								i++
						}
				}
		}
		return i
}

// 是否有变量表达式
func HasVarExpress(v string) bool {
		return expressVarReg.MatchString(v)
}

// 获取当前应用的数据库驱动
func GetDatabaseDriver() string {
		driver := beego.AppConfig.String("database")
		if driver == "" {
				return "mysql"
		}
		return driver
}

// 字符串数组
func GetKvStrArr(key string, def ...[]string) []string {
		if len(def) == 0 {
				def = append(def, []string{})
		}
		var vars = GetParsedKvStr(key)
		if vars == "" {
				return def[0]
		}
		if !strings.Contains(vars, ",") {
				return []string{vars}
		}
		arr := strings.SplitN(vars, ",", -1)
		return array.Filter(arr)
}


