package transforms

import "github.com/WebGameLinux/cms/models/enums"

func TransformGender(v interface{}) enums.Gender {
		if n, ok := v.(enums.Gender); ok {
				return n
		}
		if n, ok := v.(int); ok {
				return enums.ParseInt(n)
		}
		if str, ok := v.(string); ok {
				return enums.Parse(str)
		}
		return enums.Unknown
}

func TransformMapGender(key string, v interface{}) (string, interface{}) {
		return key, TransformGender(v)
}
